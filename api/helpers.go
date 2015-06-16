// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package api

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/dustin/go-humanize"
	"github.com/scaleway/scaleway-cli/utils"
)

// ScalewayResolvedIdentifier represents a list of matching identifier for a specifier pattern
type ScalewayResolvedIdentifier struct {
	// Identifiers holds matching identifiers
	Identifiers []ScalewayIdentifier

	// Needle is the criteria used to lookup identifiers
	Needle string
}

// ScalewayImageInterface is an interface to multiple Scaleway items
type ScalewayImageInterface struct {
	CreationDate time.Time
	Identifier   string
	Name         string
	Tag          string
	VirtualSize  float64
	Public       bool
	Type         string
}

// CreateVolumeFromHumanSize creates a volume on the API with a human readable size
func CreateVolumeFromHumanSize(api *ScalewayAPI, size string) (*string, error) {
	bytes, err := humanize.ParseBytes(size)
	if err != nil {
		return nil, err
	}

	var newVolume ScalewayVolumeDefinition
	newVolume.Name = size
	newVolume.Size = bytes
	newVolume.Type = "l_ssd"

	volumeID, err := api.PostVolume(newVolume)
	if err != nil {
		return nil, err
	}

	return &volumeID, nil
}

// fillIdentifierCache fills the cache by fetching fro the API
func fillIdentifierCache(api *ScalewayAPI) {
	log.Debugf("Filling the cache")
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		api.GetServers(true, 0)
		wg.Done()
	}()
	go func() {
		api.GetImages()
		wg.Done()
	}()
	go func() {
		api.GetSnapshots()
		wg.Done()
	}()
	go func() {
		api.GetVolumes()
		wg.Done()
	}()
	go func() {
		api.GetBootscripts()
		wg.Done()
	}()
	wg.Wait()
}

// GetIdentifier returns a an identifier if the resolved needles only match one element, else, it exists the program
func GetIdentifier(api *ScalewayAPI, needle string) *ScalewayIdentifier {
	idents := ResolveIdentifier(api, needle)

	if len(idents) == 1 {
		return &idents[0]
	}
	if len(idents) == 0 {
		log.Fatalf("No such identifier: %s", needle)
	}
	log.Errorf("Too many candidates for %s (%d)", needle, len(idents))
	for _, identifier := range idents {
		// FIXME: also print the name
		log.Infof("- %s", identifier.Identifier)
	}
	os.Exit(1)
	return nil
}

// ResolveIdentifier resolves needle provided by the user
func ResolveIdentifier(api *ScalewayAPI, needle string) []ScalewayIdentifier {
	idents := api.Cache.LookUpIdentifiers(needle)
	if len(idents) > 0 {
		return idents
	}

	fillIdentifierCache(api)

	idents = api.Cache.LookUpIdentifiers(needle)
	return idents
}

// ResolveIdentifiers resolves needles provided by the user
func ResolveIdentifiers(api *ScalewayAPI, needles []string, out chan ScalewayResolvedIdentifier) {
	// first attempt, only lookup from the cache
	var unresolved []string
	for _, needle := range needles {
		idents := api.Cache.LookUpIdentifiers(needle)
		if len(idents) == 0 {
			unresolved = append(unresolved, needle)
		} else {
			out <- ScalewayResolvedIdentifier{
				Identifiers: idents,
				Needle:      needle,
			}
		}
	}
	// fill the cache by fetching from the API and resolve missing identifiers
	if len(unresolved) > 0 {
		fillIdentifierCache(api)

		for _, needle := range unresolved {
			idents := api.Cache.LookUpIdentifiers(needle)
			out <- ScalewayResolvedIdentifier{
				Identifiers: idents,
				Needle:      needle,
			}
		}
	}
	close(out)
}

// InspectIdentifiers inspects identifiers concurrently
func InspectIdentifiers(api *ScalewayAPI, ci chan ScalewayResolvedIdentifier, cj chan interface{}) {
	var wg sync.WaitGroup
	for {
		idents, ok := <-ci
		if !ok {
			break
		}
		if len(idents.Identifiers) != 1 {
			if len(idents.Identifiers) == 0 {
				log.Errorf("Unable to resolve identifier %s", idents.Needle)
			} else {
				log.Errorf("Too many candidates for %s (%d)", idents.Needle, len(idents.Identifiers))
				for _, identifier := range idents.Identifiers {
					// FIXME: also print the name
					log.Infof("- %s", identifier.Identifier)
				}
			}
		} else {
			ident := idents.Identifiers[0]
			wg.Add(1)
			go func() {
				if ident.Type == IdentifierServer {
					server, err := api.GetServer(ident.Identifier)
					if err == nil {
						cj <- server
					}
				} else if ident.Type == IdentifierImage {
					image, err := api.GetImage(ident.Identifier)
					if err == nil {
						cj <- image
					}
				} else if ident.Type == IdentifierSnapshot {
					snap, err := api.GetSnapshot(ident.Identifier)
					if err == nil {
						cj <- snap
					}
				} else if ident.Type == IdentifierVolume {
					snap, err := api.GetVolume(ident.Identifier)
					if err == nil {
						cj <- snap
					}
				} else if ident.Type == IdentifierBootscript {
					bootscript, err := api.GetBootscript(ident.Identifier)
					if err == nil {
						cj <- bootscript
					}
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	close(cj)
}

// CreateServer creates a server using API based on typical server fields
func CreateServer(api *ScalewayAPI, imageName string, name string, bootscript string, env string, additionalVolumes string) (string, error) {
	if name == "" {
		name = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	}

	var server ScalewayServerDefinition
	server.Volumes = make(map[string]string)

	server.Tags = []string{}
	if env != "" {
		server.Tags = strings.Split(env, " ")
	}
	if additionalVolumes != "" {
		volumes := strings.Split(additionalVolumes, " ")
		for i := range volumes {
			volumeID, err := CreateVolumeFromHumanSize(api, volumes[i])
			if err != nil {
				return "", err
			}

			volumeIDx := fmt.Sprintf("%d", i+1)
			server.Volumes[volumeIDx] = *volumeID
		}
	}
	server.Name = name
	if bootscript != "" {
		bootscript := api.GetBootscriptID(bootscript)
		server.Bootscript = &bootscript
	}

	inheritingVolume := false
	_, err := humanize.ParseBytes(imageName)
	if err == nil {
		// Create a new root volume
		volumeID, err := CreateVolumeFromHumanSize(api, imageName)
		if err != nil {
			return "", err
		}
		server.Volumes["0"] = *volumeID
	} else {
		// Use an existing image
		// FIXME: handle snapshots
		inheritingVolume = true
		image := api.GetImageID(imageName, false)
		if image != "" {
			server.Image = &image
		} else {
			snapshotID := api.GetSnapshotID(imageName)
			snapshot, err := api.GetSnapshot(snapshotID)
			if err != nil {
				return "", err
			}
			server.Volumes["0"] = snapshot.BaseVolume.Identifier
		}
	}

	serverID, err := api.PostServer(server)
	if err != nil {
		return "", nil
	}

	// For inherited volumes, we prefix the name with server hostname
	if inheritingVolume {
		createdServer, err := api.GetServer(serverID)
		if err != nil {
			return "", err
		}
		currentVolume := createdServer.Volumes["0"]

		var volumePayload ScalewayVolumePutDefinition
		newName := fmt.Sprintf("%s-%s", createdServer.Hostname, currentVolume.Name)
		volumePayload.Name = &newName
		volumePayload.CreationDate = &currentVolume.CreationDate
		volumePayload.Organization = &currentVolume.Organization
		volumePayload.Server.Identifier = &currentVolume.Server.Identifier
		volumePayload.Server.Name = &currentVolume.Server.Name
		volumePayload.Identifier = &currentVolume.Identifier
		volumePayload.Size = &currentVolume.Size
		volumePayload.ModificationDate = &currentVolume.ModificationDate
		volumePayload.ExportURI = &currentVolume.ExportURI
		volumePayload.VolumeType = &currentVolume.VolumeType

		err = api.PutVolume(currentVolume.Identifier, volumePayload)
		if err != nil {
			return "", err
		}
	}

	return serverID, nil
}

// WaitForServerState asks API in a loop until a server matches a wanted state
func WaitForServerState(api *ScalewayAPI, serverID string, targetState string) (*ScalewayServer, error) {
	var server *ScalewayServer
	var err error

	for {
		server, err = api.GetServer(serverID)
		if err != nil {
			return nil, err
		}
		if server.State == targetState {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return server, nil
}

// WaitForServerReady wait for a server state to be running, then wait for the SSH port to be available
func WaitForServerReady(api *ScalewayAPI, serverID string) (*ScalewayServer, error) {
	server, err := WaitForServerState(api, serverID, "running")
	if err != nil {
		return nil, err
	}

	dest := fmt.Sprintf("%s:22", server.PublicAddress.IP)

	err = utils.WaitForTCPPortOpen(dest)
	if err != nil {
		return nil, err
	}

	return server, nil
}

// WaitForServerStopped wait for a server state to be stopped
func WaitForServerStopped(api *ScalewayAPI, serverID string) (*ScalewayServer, error) {
	server, err := WaitForServerState(api, serverID, "stopped")
	if err != nil {
		return nil, err
	}
	return server, nil
}

// ByCreationDate sorts images by CreationDate field
type ByCreationDate []ScalewayImageInterface

func (a ByCreationDate) Len() int           { return len(a) }
func (a ByCreationDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreationDate) Less(i, j int) bool { return a[j].CreationDate.Before(a[i].CreationDate) }

// StartServer start a server based on its needle, can optionaly block while server is booting
func StartServer(api *ScalewayAPI, needle string, wait bool) error {
	server := api.GetServerID(needle)

	err := api.PostServerAction(server, "poweron")
	if err != nil {
		if err.Error() != "server should be stopped" {
			return fmt.Errorf("server %s is already started: %v", server, err)
		}
	}

	if wait {
		_, err = WaitForServerReady(api, server)
		if err != nil {
			return fmt.Errorf("failed to wait for server %s to be ready, %v", needle, err)
		}
	}
	return nil
}

// StartServerOnce wraps StartServer for golang channel
func StartServerOnce(api *ScalewayAPI, needle string, wait bool, successChan chan bool, errChan chan error) {
	err := StartServer(api, needle, wait)

	if err != nil {
		errChan <- err
		return
	}

	fmt.Println(needle)
	successChan <- true
}
