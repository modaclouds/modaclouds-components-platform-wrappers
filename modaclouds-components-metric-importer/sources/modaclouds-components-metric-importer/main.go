

package main


import "errors"
import "fmt"
import "net"
import "os"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


var selfGroup = ComponentGroup ("a4f6135ee4897fe65e257a4e7b445c213676adf6")
var explorerGroup = ComponentGroup ("7e079e121d36b8cec279b3116b2ec3d9a4e36045")


type callbacks struct {
	httpIp net.IP
	httpPort uint16
	httpFqdn string
	graphiteIp net.IP
	graphitePort uint16
	graphiteFqdn string
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
	_server.Transcript.TraceInformation ("acquiring the HTTP endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketAcquire (ResourceIdentifier ("http")); _error != nil {
		return _error
	} else {
		_callbacks.httpIp = _ip_1
		_callbacks.httpPort = _port_1
		_callbacks.httpFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the HTTP endpoint: `%s:%d`;", _callbacks.httpIp.String (), _callbacks.httpPort)
	
	_server.Transcript.TraceInformation ("resolving the metric explorer line receiver endpoint...")
	if _ip_1, _port_1, _fqdn_1, _error := _server.TcpSocketResolve (explorerGroup, "modaclouds-metric-explorer:get-line-receiver-endpoint"); _error != nil {
		return _error
	} else {
		_callbacks.graphiteIp = _ip_1
		_callbacks.graphitePort = _port_1
		_callbacks.graphiteFqdn = _fqdn_1
	}
	
	_server.Transcript.TraceInformation ("  * using the metric explorer line receiver endpoint: `%s:%d`;", _callbacks.graphiteIp.String (), _callbacks.graphitePort)
	
	_server.ProcessExecutable = os.Getenv ("modaclouds_service_run")
	
	_server.ProcessEnvironment = map[string]string {
			"MODACLOUDS_METRIC_IMPORTER_ENDPOINT_IP" : _callbacks.httpIp.String (),
			"MODACLOUDS_METRIC_IMPORTER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.httpPort),
			"MODACLOUDS_METRIC_EXPLORER_LINE_RECEIVER_ENDPOINT_IP" : _callbacks.graphiteIp.String (),
			"MODACLOUDS_METRIC_EXPLORER_LINE_RECEIVER_ENDPOINT_PORT" : fmt.Sprintf ("%d", _callbacks.graphitePort),
			"modaclouds_service_identifier" : string (_server.Identifier),
			"modaclouds_service_temporary" : fmt.Sprintf ("%s/service", _server.Temporary),
	}
	_server.SelfGroup = selfGroup
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	switch _operation {
		
		case "modaclouds-metric-importer:get-http-endpoint" :
			
			_outputs = map[string]interface{} {
					"ip" : _callbacks.httpIp.String (),
					"port" : _callbacks.httpPort,
					"fqdn" : _callbacks.httpFqdn,
					"url" : fmt.Sprintf ("http://%s:%d/", _callbacks.httpFqdn, _callbacks.httpPort),
			}
		
		default :
			
			_error = errors.New ("invalid-operation")
	}
	
	return _outputs, _error
}


func main () () {
	PreMain (& callbacks {}, packageTranscript)
}


var packageTranscript = transcript.NewPackageTranscript (transcript.InformationLevel)
