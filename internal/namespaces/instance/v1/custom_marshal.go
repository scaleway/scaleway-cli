package instance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	// register external types custom human marshaler func
	// for internal types use the human.Marshaler interface (see code at the end of file)

	// Image
	human.RegisterMarshalerFunc(instance.CreateImageResponse{}, marshallNestedField("Image"))

	// IP
	human.RegisterMarshalerFunc(instance.CreateIPResponse{}, marshallNestedField("IP"))

	// Placement Group
	human.RegisterMarshalerFunc(instance.CreatePlacementGroupResponse{}, marshallNestedField("PlacementGroup"))

	// Security Group
	human.RegisterMarshalerFunc(instance.CreateSecurityGroupResponse{}, marshallNestedField("SecurityGroup"))
	human.RegisterMarshalerFunc(instance.SecurityGroupPolicy(0), human.BindAttributesMarshalFunc(securityGroupPolicyAttribute))

	// Security Group Rule
	human.RegisterMarshalerFunc(instance.CreateSecurityGroupRuleResponse{}, marshallNestedField("Rule"))
	human.RegisterMarshalerFunc(instance.SecurityGroupRuleAction(0), human.BindAttributesMarshalFunc(securityGroupRuleActionAttribute))

	// Server
	human.RegisterMarshalerFunc(instance.CreateServerResponse{}, marshallNestedField("Server"))
	human.RegisterMarshalerFunc(instance.ServerState(0), serverStateMarshallerFunc)
	human.RegisterMarshalerFunc(instance.ServerLocation{}, serverLocationMarshallerFunc)
	human.RegisterMarshalerFunc([]*instance.Server{}, serversMarshallerFunc)
	human.RegisterMarshalerFunc(instance.GetServerResponse{}, getServerResponseMarshallerFunc)
	human.RegisterMarshalerFunc(instance.Bootscript{}, bootscriptMarshallerFunc)

	// Snapshot
	human.RegisterMarshalerFunc(instance.CreateSnapshotResponse{}, marshallNestedField("Snapshot"))

	// Volume
	human.RegisterMarshalerFunc(instance.CreateVolumeResponse{}, marshallNestedField("Volume"))
	human.RegisterMarshalerFunc(instance.VolumeState(0), human.BindAttributesMarshalFunc(volumeStateAttributes))
	human.RegisterMarshalerFunc(instance.VolumeSummary{}, volumeSummaryMarshallerFunc)
	human.RegisterMarshalerFunc(map[string]*instance.Volume{}, volumeMapMarshallerFunc)
}

// serverStateMarshallerFunc marshals a instance.ServerState.
var (
	serverStateAttributes = human.Attributes{
		instance.ServerStateRunning:        color.FgGreen,
		instance.ServerStateStopped:        color.Faint,
		instance.ServerStateStoppedInPlace: color.Faint,
		instance.ServerStateStarting:       color.FgBlue,
		instance.ServerStateStopping:       color.FgBlue,
		instance.ServerStateLocked:         color.FgRed,
	}

	volumeStateAttributes = human.Attributes{
		instance.VolumeStateError:     color.FgRed,
		instance.VolumeStateAvailable: color.FgGreen,
	}

	securityGroupPolicyAttribute = human.Attributes{
		instance.SecurityGroupPolicyDrop:   color.FgRed,
		instance.SecurityGroupPolicyAccept: color.FgGreen,
	}

	securityGroupRuleActionAttribute = human.Attributes{
		instance.SecurityGroupRuleActionDrop:   color.FgRed,
		instance.SecurityGroupRuleActionAccept: color.FgGreen,
	}
)

// serverLocationMarshallerFunc marshals a instance.ServerLocation.
func serverLocationMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	location := i.(instance.ServerLocation)
	zone, err := scw.ParseZone(location.ZoneID)
	if err != nil {
		return "", err
	}
	zoneStr := fmt.Sprintf("%s", zone)
	return zoneStr, nil
}

// serverStateMarshallerFunc marshals a instance.ServerState.
func serverStateMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// The Scaleway console shows "archived" for a stopped server.
	if i.(instance.ServerState) == instance.ServerStateStopped {
		return terminal.Style("archived", color.Faint), nil
	}
	return human.BindAttributesMarshalFunc(serverStateAttributes)(i, opt)
}

// serversMarshallerFunc marshals a Server.
func serversMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// humanServerInList is the custom Server type used for list view.
	type humanServerInList struct {
		ID                string
		Name              string
		State             instance.ServerState
		Zone              scw.Zone
		PublicIP          net.IP
		PrivateIP         *string
		ImageName         string
		Tags              []string
		ModificationDate  time.Time
		CreationDate      time.Time
		ImageId           string
		Protected         bool
		Volumes           int
		SecurityGroupId   string
		SecurityGroupName string
		StateDetail       string
		Arch              instance.Arch
		PlacementGroup    *instance.PlacementGroup
	}

	servers := i.([]*instance.Server)
	humanServers := make([]*humanServerInList, 0)
	for _, server := range servers {
		var zone scw.Zone
		if server.Location != nil {
			zone_, err := scw.ParseZone(server.Location.ZoneID)
			if err != nil {
				return "", err
			}
			zone = zone_
		}
		publicIPAddress := net.IP(nil)
		if server.PublicIP != nil {
			publicIPAddress = server.PublicIP.Address
		}
		serverImageID := ""
		serverImageName := ""
		if server.Image != nil {
			serverImageID = server.Image.ID
			serverImageName = server.Image.Name
		}
		humanServers = append(humanServers, &humanServerInList{
			ID:                server.ID,
			Name:              server.Name,
			State:             server.State,
			Zone:              zone,
			ModificationDate:  server.ModificationDate,
			CreationDate:      server.CreationDate,
			ImageId:           serverImageID,
			ImageName:         serverImageName,
			Protected:         server.Protected,
			PublicIP:          publicIPAddress,
			PrivateIP:         server.PrivateIP,
			Volumes:           len(server.Volumes),
			SecurityGroupId:   server.SecurityGroup.ID,
			SecurityGroupName: server.SecurityGroup.Name,
			StateDetail:       server.StateDetail,
			Arch:              server.Arch,
			PlacementGroup:    server.PlacementGroup,
			Tags:              server.Tags,
		})
	}
	return human.Marshal(humanServers, opt)
}

// serversMarshallerFunc marshals a BootscriptID.
func bootscriptMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	bootscript := i.(instance.Bootscript)
	return bootscript.Title, nil
}

// serversMarshallerFunc marshals a VolumeSummary.
func volumeSummaryMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumeSummary := i.(instance.VolumeSummary)
	return human.Marshal(volumeSummary.ID, opt)
}

// volumeMapMarshallerFunc returns the length of the map.
func volumeMapMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumes := i.(map[string]*instance.Volume)
	return fmt.Sprintf("%v", len(volumes)), nil
}

func getServerResponseMarshallerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	serverResponse := i.(instance.GetServerResponse)

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "server",
			Title:     "Server",
		},
		{
			FieldName: "server.image",
			Title:     "Server Image",
		}, {
			FieldName: "server.allowed-actions",
		}, {
			FieldName: "volumes",
			Title:     "Volumes",
		},
	}

	customServer := &struct {
		Server  *instance.Server
		Volumes []*instance.Volume
	}{
		serverResponse.Server,
		orderVolumes(serverResponse.Server.Volumes),
	}

	str, err := human.Marshal(customServer, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

// orderVolumes return an ordered slice based on the volume map key "0", "1", "2",...
func orderVolumes(v map[string]*instance.Volume) []*instance.Volume {
	indexes := []string(nil)
	for index := range v {
		indexes = append(indexes, index)
	}
	sort.Strings(indexes)
	var orderedVolumes []*instance.Volume
	for _, index := range indexes {
		orderedVolumes = append(orderedVolumes, v[index])
	}
	return orderedVolumes
}

// marshallNestedField will marshal only the given field of a struct.
func marshallNestedField(nestedKey string) human.MarshalerFunc {
	return func(i interface{}, opt *human.MarshalOpt) (s string, err error) {
		if reflect.TypeOf(i).Kind() != reflect.Struct {
			return "", fmt.Errorf("%T must be a struct", i)
		}
		nestedValue := reflect.ValueOf(i).FieldByName(nestedKey)
		return human.Marshal(nestedValue.Interface(), opt)
	}
}

////
// Type implementing human.Marshaler interface
////

// MarshalHuman marshals a customSecurityGroupResponse.
func (sg *customSecurityGroupResponse) MarshalHuman() (out string, err error) {
	humanSecurityGroup := struct {
		ID                    string
		Name                  string
		Description           string
		EnableDefaultSecurity bool
		OrganizationID        string
		OrganizationDefault   bool
		CreationDate          time.Time
		ModificationDate      time.Time
		Stateful              bool
	}{
		ID:                    sg.ID,
		Name:                  sg.Name,
		Description:           sg.Description,
		EnableDefaultSecurity: sg.EnableDefaultSecurity,
		OrganizationID:        sg.Organization,
		OrganizationDefault:   sg.OrganizationDefault,
		CreationDate:          sg.CreationDate,
		ModificationDate:      sg.ModificationDate,
		Stateful:              sg.Stateful,
	}

	securityGroupView, err := human.Marshal(humanSecurityGroup, nil)
	if err != nil {
		return "", err
	}
	securityGroupView = terminal.Style("Security Group:\n", color.Bold) + securityGroupView

	type humanRule struct {
		ID       string
		Protocol instance.SecurityGroupRuleProtocol
		Action   instance.SecurityGroupRuleAction
		IPRange  string
		Dest     string
	}

	toHumanRule := func(rule *instance.SecurityGroupRule) *humanRule {
		dest := "ALL"
		if rule.DestPortFrom != nil {
			dest = strconv.Itoa(int(*rule.DestPortFrom))
		}
		if rule.DestPortTo != nil {
			dest += "-" + strconv.Itoa(int(*rule.DestPortTo))
		}
		return &humanRule{
			ID:       rule.ID,
			Protocol: rule.Protocol,
			Action:   rule.Action,
			IPRange:  rule.IPRange.String(),
			Dest:     dest,
		}
	}

	inboundRules := []*humanRule(nil)
	outboundRules := []*humanRule(nil)
	for _, rule := range sg.Rules {
		switch rule.Direction {
		case instance.SecurityGroupRuleDirectionInbound:
			inboundRules = append(inboundRules, toHumanRule(rule))
		case instance.SecurityGroupRuleDirectionOutbound:
			outboundRules = append(outboundRules, toHumanRule(rule))
		default:
			logger.Warningf("invalid security group rule direction: %v", rule.Direction)
		}
	}

	defaultInboundPolicy, err := human.Marshal(sg.InboundDefaultPolicy, nil)
	if err != nil {
		return "", err
	}

	defaultOutboundPolicy, err := human.Marshal(sg.OutboundDefaultPolicy, nil)
	if err != nil {
		return "", err
	}

	inboundRulesContent, err := human.Marshal(inboundRules, nil)
	if err != nil {
		return "", err
	}
	inboundRulesView := terminal.Style(fmt.Sprintf("Inbound Rules (default policy %s):\n", defaultInboundPolicy), color.Bold) + inboundRulesContent

	outboundRulesContent, err := human.Marshal(outboundRules, nil)
	if err != nil {
		return "", err
	}
	outboundRulesView := terminal.Style(fmt.Sprintf("Outbound Rules (default policy %s):\n", defaultOutboundPolicy), color.Bold) + outboundRulesContent

	return strings.Join([]string{securityGroupView, inboundRulesView, outboundRulesView}, "\n\n"), nil
}
