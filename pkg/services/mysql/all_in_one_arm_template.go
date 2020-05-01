package mysql

// nolint: lll
var allInOneARMTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		},
		"administratorLoginPassword": {
			"type": "securestring"
		}
	},
	"variables": {
		"DBforMySQLapiVersion": "2017-12-01"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforMySQLapiVersion')]",
			"kind": "",
			"location": "{{.location}}",
			"name": "{{ .serverName }}",
			"properties": {
				"version": "{{.version}}",
				"administratorLogin": "{{ .administratorLogin }}",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"storageProfile": {
					"storageMB": {{.storage}},
					{{ if .geoRedundantBackup }}
					"geoRedundantBackup": "Enabled",
					{{ end }}
					"backupRetentionDays": {{.backupRetention}}
				},
				"sslEnforcement": "{{ .sslEnforcement }}"
			},
			"sku": {
				"name": "{{.sku}}",
				"tier": "{{.tier}}",
				"capacity": "{{.cores}}",
				"size": "{{.storage}}",
				"family": "{{.family}}"
			},
			"type": "Microsoft.DBforMySQL/servers",
			"tags": "[parameters('tags')]",
			"resources": [
				{{ $root := . }}
				{{$firewallRulesCount := sub (len .firewallRules) 1 }}
				{{$virtualNetworkRulesCount := sub (len .virtualNetworkRules) 1 }}
				{{range $i, $rule := .firewallRules}}
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforMySQLapiVersion')]",
					"dependsOn": [
						"Microsoft.DBforMySQL/servers/{{ $.serverName }}"
					],
					"location": "{{$root.location}}",
					"name": "{{$rule.name}}",
					"properties": {
						"startIpAddress": "{{$rule.startIPAddress}}",
						"endIpAddress": "{{$rule.endIPAddress}}"
					}
				},
				{{end}}	
				{{range $i, $rule := .virtualNetworkRules}}
				{
					"type": "virtualNetworkRules",
					"apiVersion": "[variables('DBforMySQLapiVersion')]",
					"dependsOn": [
						"Microsoft.DBforMySQL/servers/{{ $.serverName }}"
					],
					"location": "{{$root.location}}",
					"name": "{{$rule.name}}",
					"properties": {
						"virtualNetworkSubnetId": "[resourceId('{{ $root.virtualNetworkResourceGroup }}','Microsoft.Network/virtualNetworks/subnets', '{{ $root.virtualNetworkName }}', '{{ $rule.subnetName }}')]"
					}
				},
				{{end}}
				{
					"apiVersion": "[variables('DBforMySQLapiVersion')]",
					"name": "{{ .databaseName }}",
					"type": "databases",
					"location": "{{$root.location}}",
					"dependsOn": [
						{{range $.firewallRules}}
						"Microsoft.DBforMySQL/servers/{{ $.serverName }}/firewallrules/{{.name}}",
						{{end}}
						"Microsoft.DBforMySQL/servers/{{ $.serverName }}"
					],
					"properties": {}
				}
			]
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference('{{ .serverName }}').fullyQualifiedDomainName]"
		}
	}
}
`)
