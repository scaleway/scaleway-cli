package k8s

import (
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/ghodss/yaml"
	api "github.com/kubernetes-client/go-base/config/api"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

var (
	existingKubeconfig = `apiVersion: v1
kind: Config
current-context: test@test
clusters:
- name: test
  cluster:
    server: https://test:6443
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJd01ETXdOREV4TVRVd01Gb1hEVE13TURNd05ERXhNVFV3TUZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS3UvCkxIWXlNUTdFWmRaWmFaRUFsMkJaZDJLVmhPNjR0cGdJN0NUcG56NGRMUzZpdmRQSVJMbGdGdzBBV3ByNkFQeGsKeWp4K09mMkxxUVR4TkdDTFFaOEN2LzF4K3ZSdlM5YXJ3VnVIZ1pML2ZkTTE5Z2oyRWMvZUY3WDhLMEVadVNYMQpHb1ZickFOSDcyb2NkeHdvajZZc29nZHlSRHpXbkRjN2pvVjlsZlZjMlF1UDdHK2FGOEZBSTYrMUZpTDlidFE3Cnl3ZmFDWnE4M1dKcTFpOStqKzRsRUFDaFA1alZ2LzBLenVOMnRud3VaV3ljM3pRbHJHSG1xbGQ4WCtjVEJrMGgKN3hYSTV3U3VPTzMvVzBpcFQ5QVdWZjBoWVpnNVY5VWlKTjVBREZXem15eFQ4bTlVZXR3bTNhVlJtZ3IrNU44YQpGNkdWTWd2TEFteVo1bUpoR0RrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFFL1JMZU9vbXM3TlF4OUgvNjJIcm9VUTJ5TFcKMXlpaVBVaFBEUUJEZ1p5T251TWN2Sk84SFQrcFZ3OWVPazUxeTNUT01LYmJwQWhxc0I3dEpLQVRWcWJ5Z2tHNgpzWUQrS3kwckNITzdEMTM5RXQ1eEI2dGtDQTJxSlVRMkFFb0kyNy9na3Z3K05WektUU3VSVUF0QXZXQ2diNnVZCmRUa2NpaGNSemhyV1dTV05RN3M2RSt0Y3NJS01ibUJVT0NsMVh6MERPM3RsbzlPLzRMUkE4b3RXUDEweGVLcmsKVWxRVUsvSXFpL2FTN1c4Nm4vSWU1dFVvNTJRZk5NMk0xZkVtUVA2eThYMUZWbDh2cC9WcXgxMWszeS9ZSEcvVgp5N0pIY1hhdVJLSTdieGdPdWdTZXAyZDBsaFl4Smw2NUNQT0NsbGF3Y3U3MERzMzRNS2kzWGtDZTIwST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
contexts:
- name: test@test
  context:
    cluster: test
    user: test
users:
- name: test
  user:
    token: qotGxuOfD74ajgWir18tMMPicxLIizKg3nt5PKHGbsbNDbGfqNojIdXI
    username: test
    password: test
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvekNDQWVlZ0F3SUJBZ0lJZERQak80Umphdzh3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TURBek1qTXdPRFEyTkRoYUZ3MHlNVEF6TWpNd09EUTJORGhhTUVFeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1TWXdKQVlEVlFRREV4MXJkV0psTFdGd2FYTmxjblpsCmNpMXJkV0psYkdWMExXTnNhV1Z1ZERDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUIKQU54VjByQ0lBemNsdXIyV1VNb3NqOW1LQmlkclYzcnB5RmNwdnltMmtFVjZaOVo2TTBSRXpyTHo1c3BaWndCTwo1bHZrbEdzL2RJVndFK2pBd2tNWWNRRWlOaTQ2bHU4UFNSei9HVTFkek5mOEF2TXpnRWZER0xUY2x3eUs4di9kCklLenhTUnVOUFFseDZoTUw1bFpDeVBBZ3hqejNEdDZGWmUxUnVUdURWTUhnOWZIaHNwOFZTYnVCbWFYTTU2T0IKLzNZQXJLMXZOTlY0enRlQ3libFZnVUd3QUdKQ09zTlE0d0l4R0xSdjN5TVhtK3V3YVpGeTFxSEh6ZlpXclRpQQpKQ2lQNFVCbDV3bnUzeEhNaFZaemI0RnNCLzBmVEl1WHQ0ZjQ5L201KzdpM01vMEdrMjJNMjAvQldzNURZVmo1CnptSVVxcU9kK09UekdkcjgvcTRsdnQ4Q0F3RUFBYU1uTUNVd0RnWURWUjBQQVFIL0JBUURBZ1dnTUJNR0ExVWQKSlFRTU1Bb0dDQ3NHQVFVRkJ3TUNNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUNuVFFCWlhsbm1aVWpDNEdscwpKdTZWWEYxN040ckpzWkNVOVR3SEhETlpZNGo5YlZHbG8wZzAva3JTajBId3hNNVU1NXl2WUJDaWNpZzVkSS96Cnd2ZENUQm5FQWIxRWtuZVR1ZkVPYzFtNzBtSzg0dnd3WWZtRVNkY1NXMmJieHBuUFNpak5BdnlTekZTTmZZZDEKMy9FZlRlQjQ0VFNGRGZQVk83YnpKYXBpYVJCTlZocVJQSncwc0lJWGM1Q0hiQzFEMHU5Mk4zRnhCa3JKcFN2UAp1QXBQT2dyNUgwUk5rOEk2TTBjd0FBc1RqdUkxd2Z4MjhJU0FWcmZLUjU4d1Eza1NsZzlUTTQrN01VMFA4eUZHClJXRkIrVFZiMTExYTRDc2RSbWMzQnZtcnFEbjZ2Ny9LOTJ4c0hNeDdBd3FObk1XUDQ4QStoVFNFVFh3U1Btb3cKL040RAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBM0ZYU3NJZ0ROeVc2dlpaUXlpeVAyWW9HSjJ0WGV1bklWeW0vS2JhUVJYcG4xbm96ClJFVE9zdlBteWxsbkFFN21XK1NVYXo5MGhYQVQ2TURDUXhoeEFTSTJManFXN3c5SkhQOFpUVjNNMS93Qzh6T0EKUjhNWXROeVhESXJ5LzkwZ3JQRkpHNDA5Q1hIcUV3dm1Wa0xJOENER1BQY08zb1ZsN1ZHNU80TlV3ZUQxOGVHeQpueFZKdTRHWnBjem5vNEgvZGdDc3JXODAxWGpPMTRMSnVWV0JRYkFBWWtJNncxRGpBakVZdEcvZkl4ZWI2N0JwCmtYTFdvY2ZOOWxhdE9JQWtLSS9oUUdYbkNlN2ZFY3lGVm5OdmdXd0gvUjlNaTVlM2gvajMrYm43dUxjeWpRYVQKYll6YlQ4RmF6a05oV1BuT1loU3FvNTM0NVBNWjJ2eityaVcrM3dJREFRQUJBb0lCQVFESDRsdldwaTAwbEZmSwpzbGpzY0d5M2p3MXlLV0VkTW9UNi9mWmNJekRTdHU4SWxhZDRvV3RhMFFWb1FKNittdFZFUENPZy85bjNTK3ZqCjFTcm1yMytrNWFKOVljMlhaaWlQMDZUaW1OdkNmTzg0TGxxTHY2UGtQOUlRSU9XOTFKOVdCVGFyZGdBUFYzWmcKZlFVaThFZFdBSVdXdlJLU01EWjlpd3dkdjFEZTZFUmt4Z0Y2R0NTSXQ2Ri80RS81Uk1VbkJObU1ycjZHWHR5NgorK0cxWExCcWxRdExYVm1yRDAyVW05Y1Yyb2QwOEczTzdUM3VqUWl6ZjR6emx4LzVWWGk0ZTFkVEViY05PRU53Cnlwd3lSajBCdFh2TXVwUXZvdUZRM3I5UmVQL0g5dmp1Q3NiZGF3T1pGQkFDb3J0UVJxcnFodENZMERRK2tiM1AKQWV1SjNnb0JBb0dCQU56QU0rbzBqckhKWGZJVFppSWEwZkV4QkIyNHBpTE5NTmoxaHIxSEJFWjI2eXFucFg5UApTSkRIbXhWREo2UXROZFQxZ3Y5L1MwelF3S2ZzYVZ5M2VYNW9OcU5hVEZERGhPSEoxWDZZUElwREZGTUgwNGV2CnRXV1ZNd21MVU9mdmhQR3NYOW5rdFRlVmxueTlnSUZOK0dkWFRTSlgyOEVIaHE2NGg4ay9IVXFCQW9HQkFQK0UKb2pVcUN3RWJ0UTIzcGRwYjNGRnNwdXZWM0F1aENiNnNmcUlxQ1ovVDRlUXJSSWtPU2luYmlva042ZFR5MVhuNQp6cGlJTEhOQ20wYkl2cVpJZmhkdERsUlcwcGQxbmlGZ2R3c3FacjdFUlFlN29XSHZkbVRwa1NaQ2p1M04zb1NjCjRPSmUwVmxBdWdwMjRsbms2bisxb01ySjJRUjNqQkxQVGUvZ3dKbGZBb0dCQUtBcUhBQ3J6WFNWQThLbDdJNkcKSXhqNlZXQXpIdWRWTlVIVk1zT1dDVFlQQmlWV3FhOHJHUjFpbGRUaGVwdVY2ZDd2bXZKQnE2SzZPMjRiQzM4bgo1OUNkVURkSlJ1RzZXbWx3QmFUcVU5S0ZSUFBSVTlxNDA4WTJjR2RXVzRkTXM0cWRaSlkxYUg1QjNJUDVBb25PCnhwSkVOMFRadGluaGlnaXUvbVkza3NzQkFvR0JBUFZDb0ZnYmhQaUpXZDVTMnRXZnV2aEZMR3ZPbVNwb1p1d28Kc2x5QnNUOUNwOTdWVVRHbEQ3Ymh6allEcnVFQ1BicVk5NThkaGwwVUgrdHZvT0FIVVZDM0V6d05JcExUQ1BmTQptamVUZVkrKzRPdXRSQmkzTzVOZFJqL05QMWd2ZFZraEpCTGxKRmxoY2JHOXIwTE9JZkIzckdFNkloN1JpUmc4CjkvZzZhV1JOQW9HQWJSSG94d1B4MUVlRnRxVDlUdXowZWZUR3RwQTB0bkhDZTN2b2x1L201eEZ6N3BwS29HeUQKRkNPVm5jMmZ3LzQwYUFGTGdHYlFLMFBqTzFCbWZ3cjFvb09aT1hZYnExUXo3Q1cvN3A1OUFkR0VrWXFzdWZZcAp6OXlMd1dBUEdybm9jVjBVQXZ2SHcvbC9OK29NZEdpdmVTdDhRb3RHclgzdm9PTmVsWThCZDRNPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`

	testKubeconfig = api.Config{
		APIVersion: kubeconfigAPIVersion,
		Kind:       kubeconfigKind,
		Clusters: []api.NamedCluster{
			{
				Name: "test",
				Cluster: api.Cluster{
					CertificateAuthorityData: []byte(`-----BEGIN CERTIFICATE-----
MIICyDCCAbCgAwIBAgIBADANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwprdWJl
cm5ldGVzMB4XDTIwMDMwNDExMTUwMFoXDTMwMDMwNDExMTUwMFowFTETMBEGA1UE
AxMKa3ViZXJuZXRlczCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKu/
LHYyMQ7EZdZZaZEAl2BZd2KVhO64tpgI7CTpnz4dLS6ivdPIRLlgFw0AWpr6APxk
yjx+Of2LqQTxNGCLQZ8Cv/1x+vRvS9arwVuHgZL/fdM19gj2Ec/eF7X8K0EZuSX1
GoVbrANH72ocdxwoj6YsogdyRDzWnDc7joV9lfVc2QuP7G+aF8FAI6+1FiL9btQ7
ywfaCZq83WJq1i9+j+4lEAChP5jVv/0KzuN2tnwuZWyc3zQlrGHmqld8X+cTBk0h
7xXI5wSuOO3/W0ipT9AWVf0hYZg5V9UiJN5ADFWzmyxT8m9Uetwm3aVRmgr+5N8a
F6GVMgvLAmyZ5mJhGDkCAwEAAaMjMCEwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB
/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAE/RLeOoms7NQx9H/62HroUQ2yLW
1yiiPUhPDQBDgZyOnuMcvJO8HT+pVw9eOk51y3TOMKbbpAhqsB7tJKATVqbygkG6
sYD+Ky0rCHO7D139Et5xB6tkCA2qJUQ2AEoI27/gkvw+NVzKTSuRUAtAvWCgb6uY
dTkcihcRzhrWWSWNQ7s6E+tcsIKMbmBUOCl1Xz0DO3tlo9O/4LRA8otWP10xeKrk
UlQUK/Iqi/aS7W86n/Ie5tUo52QfNM2M1fEmQP6y8X1FVl8vp/Vqx11k3y/YHG/V
y7JHcXauRKI7bxgOugSep2d0lhYxJl65CPOCllawcu70Ds34MKi3XkCe20I=
-----END CERTIFICATE-----
`),
					Server: "https://test:6443",
				},
			},
		},
		Contexts: []api.NamedContext{
			{
				Name: "test@test",
				Context: api.Context{
					Cluster:  "test",
					AuthInfo: "test",
				},
			},
		},
		AuthInfos: []api.NamedAuthInfo{
			{
				Name: "test",
				AuthInfo: api.AuthInfo{
					Token:                 "qotGxuOfD74ajgWir18tMMPicxLIizKg3nt5PKHGbsbNDbGfqNojIdXI",
					Username:              "test",
					Password:              "test",
					ClientCertificateData: []byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvekNDQWVlZ0F3SUJBZ0lJZERQak80Umphdzh3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TURBek1qTXdPRFEyTkRoYUZ3MHlNVEF6TWpNd09EUTJORGhhTUVFeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1TWXdKQVlEVlFRREV4MXJkV0psTFdGd2FYTmxjblpsCmNpMXJkV0psYkdWMExXTnNhV1Z1ZERDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUIKQU54VjByQ0lBemNsdXIyV1VNb3NqOW1LQmlkclYzcnB5RmNwdnltMmtFVjZaOVo2TTBSRXpyTHo1c3BaWndCTwo1bHZrbEdzL2RJVndFK2pBd2tNWWNRRWlOaTQ2bHU4UFNSei9HVTFkek5mOEF2TXpnRWZER0xUY2x3eUs4di9kCklLenhTUnVOUFFseDZoTUw1bFpDeVBBZ3hqejNEdDZGWmUxUnVUdURWTUhnOWZIaHNwOFZTYnVCbWFYTTU2T0IKLzNZQXJLMXZOTlY0enRlQ3libFZnVUd3QUdKQ09zTlE0d0l4R0xSdjN5TVhtK3V3YVpGeTFxSEh6ZlpXclRpQQpKQ2lQNFVCbDV3bnUzeEhNaFZaemI0RnNCLzBmVEl1WHQ0ZjQ5L201KzdpM01vMEdrMjJNMjAvQldzNURZVmo1CnptSVVxcU9kK09UekdkcjgvcTRsdnQ4Q0F3RUFBYU1uTUNVd0RnWURWUjBQQVFIL0JBUURBZ1dnTUJNR0ExVWQKSlFRTU1Bb0dDQ3NHQVFVRkJ3TUNNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUNuVFFCWlhsbm1aVWpDNEdscwpKdTZWWEYxN040ckpzWkNVOVR3SEhETlpZNGo5YlZHbG8wZzAva3JTajBId3hNNVU1NXl2WUJDaWNpZzVkSS96Cnd2ZENUQm5FQWIxRWtuZVR1ZkVPYzFtNzBtSzg0dnd3WWZtRVNkY1NXMmJieHBuUFNpak5BdnlTekZTTmZZZDEKMy9FZlRlQjQ0VFNGRGZQVk83YnpKYXBpYVJCTlZocVJQSncwc0lJWGM1Q0hiQzFEMHU5Mk4zRnhCa3JKcFN2UAp1QXBQT2dyNUgwUk5rOEk2TTBjd0FBc1RqdUkxd2Z4MjhJU0FWcmZLUjU4d1Eza1NsZzlUTTQrN01VMFA4eUZHClJXRkIrVFZiMTExYTRDc2RSbWMzQnZtcnFEbjZ2Ny9LOTJ4c0hNeDdBd3FObk1XUDQ4QStoVFNFVFh3U1Btb3cKL040RAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
					ClientKeyData:         []byte("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBM0ZYU3NJZ0ROeVc2dlpaUXlpeVAyWW9HSjJ0WGV1bklWeW0vS2JhUVJYcG4xbm96ClJFVE9zdlBteWxsbkFFN21XK1NVYXo5MGhYQVQ2TURDUXhoeEFTSTJManFXN3c5SkhQOFpUVjNNMS93Qzh6T0EKUjhNWXROeVhESXJ5LzkwZ3JQRkpHNDA5Q1hIcUV3dm1Wa0xJOENER1BQY08zb1ZsN1ZHNU80TlV3ZUQxOGVHeQpueFZKdTRHWnBjem5vNEgvZGdDc3JXODAxWGpPMTRMSnVWV0JRYkFBWWtJNncxRGpBakVZdEcvZkl4ZWI2N0JwCmtYTFdvY2ZOOWxhdE9JQWtLSS9oUUdYbkNlN2ZFY3lGVm5OdmdXd0gvUjlNaTVlM2gvajMrYm43dUxjeWpRYVQKYll6YlQ4RmF6a05oV1BuT1loU3FvNTM0NVBNWjJ2eityaVcrM3dJREFRQUJBb0lCQVFESDRsdldwaTAwbEZmSwpzbGpzY0d5M2p3MXlLV0VkTW9UNi9mWmNJekRTdHU4SWxhZDRvV3RhMFFWb1FKNittdFZFUENPZy85bjNTK3ZqCjFTcm1yMytrNWFKOVljMlhaaWlQMDZUaW1OdkNmTzg0TGxxTHY2UGtQOUlRSU9XOTFKOVdCVGFyZGdBUFYzWmcKZlFVaThFZFdBSVdXdlJLU01EWjlpd3dkdjFEZTZFUmt4Z0Y2R0NTSXQ2Ri80RS81Uk1VbkJObU1ycjZHWHR5NgorK0cxWExCcWxRdExYVm1yRDAyVW05Y1Yyb2QwOEczTzdUM3VqUWl6ZjR6emx4LzVWWGk0ZTFkVEViY05PRU53Cnlwd3lSajBCdFh2TXVwUXZvdUZRM3I5UmVQL0g5dmp1Q3NiZGF3T1pGQkFDb3J0UVJxcnFodENZMERRK2tiM1AKQWV1SjNnb0JBb0dCQU56QU0rbzBqckhKWGZJVFppSWEwZkV4QkIyNHBpTE5NTmoxaHIxSEJFWjI2eXFucFg5UApTSkRIbXhWREo2UXROZFQxZ3Y5L1MwelF3S2ZzYVZ5M2VYNW9OcU5hVEZERGhPSEoxWDZZUElwREZGTUgwNGV2CnRXV1ZNd21MVU9mdmhQR3NYOW5rdFRlVmxueTlnSUZOK0dkWFRTSlgyOEVIaHE2NGg4ay9IVXFCQW9HQkFQK0UKb2pVcUN3RWJ0UTIzcGRwYjNGRnNwdXZWM0F1aENiNnNmcUlxQ1ovVDRlUXJSSWtPU2luYmlva042ZFR5MVhuNQp6cGlJTEhOQ20wYkl2cVpJZmhkdERsUlcwcGQxbmlGZ2R3c3FacjdFUlFlN29XSHZkbVRwa1NaQ2p1M04zb1NjCjRPSmUwVmxBdWdwMjRsbms2bisxb01ySjJRUjNqQkxQVGUvZ3dKbGZBb0dCQUtBcUhBQ3J6WFNWQThLbDdJNkcKSXhqNlZXQXpIdWRWTlVIVk1zT1dDVFlQQmlWV3FhOHJHUjFpbGRUaGVwdVY2ZDd2bXZKQnE2SzZPMjRiQzM4bgo1OUNkVURkSlJ1RzZXbWx3QmFUcVU5S0ZSUFBSVTlxNDA4WTJjR2RXVzRkTXM0cWRaSlkxYUg1QjNJUDVBb25PCnhwSkVOMFRadGluaGlnaXUvbVkza3NzQkFvR0JBUFZDb0ZnYmhQaUpXZDVTMnRXZnV2aEZMR3ZPbVNwb1p1d28Kc2x5QnNUOUNwOTdWVVRHbEQ3Ymh6allEcnVFQ1BicVk5NThkaGwwVUgrdHZvT0FIVVZDM0V6d05JcExUQ1BmTQptamVUZVkrKzRPdXRSQmkzTzVOZFJqL05QMWd2ZFZraEpCTGxKRmxoY2JHOXIwTE9JZkIzckdFNkloN1JpUmc4CjkvZzZhV1JOQW9HQWJSSG94d1B4MUVlRnRxVDlUdXowZWZUR3RwQTB0bkhDZTN2b2x1L201eEZ6N3BwS29HeUQKRkNPVm5jMmZ3LzQwYUFGTGdHYlFLMFBqTzFCbWZ3cjFvb09aT1hZYnExUXo3Q1cvN3A1OUFkR0VrWXFzdWZZcAp6OXlMd1dBUEdybm9jVjBVQXZ2SHcvbC9OK29NZEdpdmVTdDhRb3RHclgzdm9PTmVsWThCZDRNPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="),
				},
			},
		},
	}
)

// testIfKubeconfigInFile checks if the given kubeconfig is in the given file
// it test if the user, cluster and context of the kubeconfig file are in the given file
func testIfKubeconfigInFile(t *testing.T, filePath string, suffix string, kubeconfig api.Config) {
	kubeconfigBytes, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	var existingKubeconfig api.Config
	err = yaml.Unmarshal(kubeconfigBytes, &existingKubeconfig)
	assert.Nil(t, err)

	found := false
	for _, cluster := range existingKubeconfig.Clusters {
		if cluster.Name == kubeconfig.Clusters[0].Name+suffix {
			//t.Log(string(cluster.Cluster.CertificateAuthorityData))
			// t.Log(string(kubeconfig.Clusters[0].Cluster.CertificateAuthorityData))
			// panic(string(cluster.Cluster.CertificateAuthorityData))
			//panic(string(kubeconfig.Clusters[0].Cluster.CertificateAuthorityData))
			assert.Equal(t, string(kubeconfig.Clusters[0].Cluster.CertificateAuthorityData), string(cluster.Cluster.CertificateAuthorityData))
			assert.Equal(t, kubeconfig.Clusters[0].Cluster.Server, cluster.Cluster.Server)
			found = true
			break
		}
	}
	assert.True(t, found, "cluster not found in kubeconfig for cluster with suffix %s", suffix)

	found = false
	for _, context := range existingKubeconfig.Contexts {
		if context.Name == kubeconfig.Contexts[0].Name+suffix {
			assert.Equal(t, kubeconfig.Contexts[0].Context.Cluster+suffix, context.Context.Cluster)
			assert.Equal(t, kubeconfig.Contexts[0].Context.AuthInfo+suffix, context.Context.AuthInfo)
			assert.Equal(t, kubeconfig.Contexts[0].Context.Namespace, context.Context.Namespace)
			found = true
			break
		}
	}
	assert.True(t, found, "context not found in kubeconfig for cluster with suffix %s", suffix)

	found = false
	for _, user := range existingKubeconfig.AuthInfos {
		if user.Name == kubeconfig.AuthInfos[0].Name+suffix {
			assert.Equal(t, kubeconfig.AuthInfos[0].AuthInfo.Username, user.AuthInfo.Username)
			found = true
			break
		}
	}
	assert.True(t, found, "user not found in kubeconfig with suffix %s", suffix)
}

func Test_InstallKubeconfig(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfig("install-kubeconfig-simple", "Cluster", "Kubeconfig", kapsuleVersion),
		Cmd:        "scw k8s kubeconfig install {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(api.Config))
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "cli-test"),
		},
	}))

	t.Run("merge", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createClusterAndWaitAndKubeconfigAndPopulateFile("install-kubeconfig-merge", "Cluster", "Kubeconfig", kapsuleVersion, path.Join(os.TempDir(), "cli-merge-test"), []byte(existingKubeconfig)),
		Cmd:        "scw k8s kubeconfig install {{ .Cluster.ID }}",
		Check: core.TestCheckCombine(
			// no golden tests since it's os specific
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-merge-test"), "-"+ctx.Meta["Cluster"].(*k8s.Cluster).ID, ctx.Meta["Kubeconfig"].(api.Config))
				testIfKubeconfigInFile(t, path.Join(os.TempDir(), "cli-merge-test"), "", testKubeconfig)
			},
			core.TestCheckExitCode(0),
		),
		AfterFunc: deleteCluster("Cluster"),
		OverrideEnv: map[string]string{
			"KUBECONFIG": path.Join(os.TempDir(), "cli-merge-test"),
		},
	}))
}
