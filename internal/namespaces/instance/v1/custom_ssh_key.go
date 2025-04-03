package instance

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"os/exec"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	sshKeyTagPrefix = "AUTHORIZED_KEY="
)

type SSHKeyFormat struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Key  string `json:"key"`
}

func marshalSSHKeys(i interface{}, opts *human.MarshalOpt) (string, error) {
	// Custom type to avoid recursion when marshaling
	type humanKey struct {
		Name string
		Type string
		Key  string
	}

	keys := i.([]*SSHKeyFormat)
	humanKeys := make([]*humanKey, len(keys))
	for i, key := range keys {
		hKey := &humanKey{
			Name: key.Name,
			Type: key.Type,
			Key:  key.Key,
		}
		if len(hKey.Key) > 16 {
			hKey.Key = hKey.Key[:16]
		}

		humanKeys[i] = hKey
	}

	return human.Marshal(humanKeys, opts)
}

func instanceSSH() *core.Command {
	return &core.Command{
		Short: `SSH Utilities`,
		Long: `Command utilities around server SSH
- Manage keys per server
- Generate ssh config`,
		Namespace: "instance",
		Resource:  "ssh",
	}
}

type sshAddKeyRequest struct {
	Zone      scw.Zone
	ServerID  string
	PublicKey string
}

func sshAddKeyCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "ssh",
		Verb:      "add-key",
		Groups:    []string{"utility"},
		Short:     "Add a public key to a server",
		Long: `Key will be added to server's tags and added to root user on next restart.
Key is expected in openssh format "(format) (key) (comment)".
The comment will be used as key name or generated
Lookup /root/.ssh/authorized_keys on your server for more information`,
		ArgsType: reflect.TypeOf(sshAddKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "server-id",
				Short: "Server to add your key to",
			},
			{
				Name:  "public-key",
				Short: "Public key you want to add to your server",
			},
			core.ZoneArgSpec(((*instance.API)(nil)).Zones()...),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*sshAddKeyRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to fetch server: %w", err)
			}

			formattedKey := FormatSSHKeyToTag(args.PublicKey)

			for i, tag := range server.Server.Tags {
				if tag == formattedKey {
					return nil, fmt.Errorf("key already exists (tags.%d)", i)
				}
			}

			tags := server.Server.Tags
			tags = append(tags, formattedKey)

			_, err = api.UpdateServer(&instance.UpdateServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
				Tags:     &tags,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update server: %w", err)
			}

			return &core.SuccessResult{
				Resource: "ssh",
				Verb:     "add-key",
			}, nil
		},
	}
}

type sshFetchKeysRequest struct {
	Zone      scw.Zone
	ProjectID string
	Username  string
}

func sshFetchKeysCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "ssh",
		Verb:      "fetch-keys",
		Groups:    []string{"utility"},
		Short:     "Fetch SSH keys from the console and install them on multiple servers",
		Long: `Keys registered via the Scaleway Console will be propagated to the selected servers.
The command 'ssh <server-ip> -t -l <username> scw-fetch-ssh-keys --upgrade' will be run on the servers matching the zone and project filters.
Keep in mind that you need to be able to connect to your server with another key than the one you want to add.
Keep in mind that SSH keys are scoped by project.`,
		ArgsType: reflect.TypeOf(sshFetchKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "project-id",
				Short:    "Fetch the keys on all servers in the given Project",
				Required: false,
			},
			{
				Name:    "username",
				Short:   "Username used for the SSH connection",
				Default: core.DefaultValueSetter("root"),
			},
			core.ZoneArgSpec(((*instance.API)(nil)).Zones()...),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*sshFetchKeysRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			listServersRequest := &instance.ListServersRequest{
				Zone: args.Zone,
			}
			if args.ProjectID != "" {
				listServersRequest.Project = &args.ProjectID
			}
			servers, err := api.ListServers(
				listServersRequest,
				scw.WithAllPages(),
				scw.WithContext(ctx),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch servers: %w", err)
			}

			for i, server := range servers.Servers {
				msg := fmt.Sprintf("Loading SSH keys on server %q", server.Name)
				if i == 0 {
					_, _ = interactive.Println(">", msg)
				} else {
					_, _ = interactive.Println("\n>", msg)
				}

				if server.State != instance.ServerStateRunning {
					_, _ = interactive.Printf("Failed: server %q is not running", server.Name)

					continue
				}

				if len(server.PublicIPs) == 0 {
					_, _ = interactive.Printf("Failed: server %q has no public IP", server.Name)

					continue
				}

				for _, publicIP := range server.PublicIPs {
					sshArgs := []string{
						publicIP.Address.String(),
						"-t",
						"-l", args.Username,
						"/usr/sbin/scw-fetch-ssh-keys",
						"--upgrade",
					}

					sshCmd := exec.Command("ssh", sshArgs...)

					_, _ = interactive.Println(sshCmd)

					exitCode, err := core.ExecCmd(ctx, sshCmd)
					if err != nil || exitCode != 0 {
						if err != nil {
							_, _ = interactive.Println("Failed:", err)
						} else {
							_, _ = interactive.Println("Failed: ssh command failed with exit code", exitCode)
						}
					} else {
						_, _ = interactive.Println("Success")
					}
				}
			}

			return &core.SuccessResult{Empty: true}, nil
		},
	}
}

type sshListKeysRequest struct {
	Zone     scw.Zone
	ServerID string
}

func sshListKeysCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "ssh",
		Verb:      "list-keys",
		Groups:    []string{"utility"},
		Short:     "List manually added public keys",
		Long: `List only keys added manually to a server using tags.
The key comment is used as key name or generated
Lookup /root/.ssh/authorized_keys on your server for more information`,
		ArgsType: reflect.TypeOf(sshListKeysRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      "Server which keys are to be listed",
				Positional: true,
				Required:   true,
			},
			core.ZoneArgSpec(((*instance.API)(nil)).Zones()...),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*sshListKeysRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to fetch server: %w", err)
			}

			keys := []*SSHKeyFormat(nil)

			for _, tag := range server.Server.Tags {
				if !isSSHKeyTag(tag) {
					continue
				}
				keys = append(keys, expandSSHKeyTag(tag))
			}

			return keys, nil
		},
	}
}

type sshRemoveKeyRequest struct {
	Zone      scw.Zone
	ServerID  string
	Name      string
	PublicKey string
}

func sshRemoveKeyCommand() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "ssh",
		Verb:      "remove-key",
		Groups:    []string{"utility"},
		Short:     "Remove a manually added public key from a server",
		Long: `Key will be remove from server's tags and removed from root user on next restart.
Keys are identified by their comment as in openssh format.
Lookup /root/.ssh/authorized_keys on your server for more information`,
		ArgsType: reflect.TypeOf(sshRemoveKeyRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "server-id",
				Short:    "Server to add your key to",
				Required: true,
			},
			{
				Name:       "name",
				Short:      "Name of the key you want to remove, has to be the key comment or the index",
				OneOfGroup: "identifier",
			},
			{
				Name:       "public-key",
				Short:      "Public key you want to remove",
				OneOfGroup: "identifier",
			},
			core.ZoneArgSpec(((*instance.API)(nil)).Zones()...),
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*sshRemoveKeyRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to fetch server: %w", err)
			}

			newTags := make([]string, 0, len(server.Server.Tags))
			removedKeys := []*SSHKeyFormat(nil)
			for _, tag := range server.Server.Tags {
				if !isSSHKeyTag(tag) {
					continue
				}
				key := expandSSHKeyTag(tag)
				if key.Name == args.Name || tag == FormatSSHKeyToTag(args.PublicKey) {
					removedKeys = append(removedKeys, expandSSHKeyTag(tag))
				} else {
					newTags = append(newTags, tag)
				}
			}

			if len(removedKeys) == 0 {
				return nil, errors.New("no key found with given filters")
			}

			_, err = api.UpdateServer(&instance.UpdateServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
				Tags:     &newTags,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update server: %w", err)
			}

			msg := "Removed keys"
			removedKeysMsg, err := human.Marshal(removedKeys, nil)
			if err == nil {
				msg = msg + "\n" + removedKeysMsg
			}

			return &core.SuccessResult{
				Resource: "ssh",
				Verb:     "remove-key",
				Message:  msg,
			}, nil
		},
	}
}

func FormatSSHKeyToTag(publicKey string) string {
	return "AUTHORIZED_KEY=" + strings.ReplaceAll(publicKey, " ", "_")
}

func isSSHKeyTag(tag string) bool {
	return strings.HasPrefix(tag, sshKeyTagPrefix)
}

func expandSSHKeyTag(tag string) *SSHKeyFormat {
	wholeKey := strings.TrimPrefix(tag, sshKeyTagPrefix)
	// key should have underscores in place of spaces
	wholeKey = strings.ReplaceAll(wholeKey, "_", " ")
	elems := strings.Split(wholeKey, " ")

	sshKey := &SSHKeyFormat{}

	if len(elems) > 0 {
		sshKey.Type = elems[0]
	}
	if len(elems) > 1 {
		sshKey.Key = elems[1]
	}
	if len(elems) > 2 {
		sshKey.Name = elems[2]
	} else {
		sshKey.Name = generateSSHKeyName(strings.Join(elems, ""))
	}

	return sshKey
}

func generateSSHKeyName(key string) string {
	sumBuffer := make([]byte, 4)
	sum := crc32.ChecksumIEEE([]byte(key))
	binary.BigEndian.PutUint32(sumBuffer, sum)

	return base64.RawStdEncoding.EncodeToString(sumBuffer)
}
