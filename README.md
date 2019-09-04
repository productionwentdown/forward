
# forward

A simple TCP proxy. Currently used in [AppVenture](https://appventure.nushigh.edu.sg/)'s internal server to port forward from host to a Hyper-V VM. 

## Usage 

```bash
$ ./forward -help
Usage of ./forward:
  -connect string
    	forward to address
  -listen string
    	listen on address (default ":8000")
  -ssh string
    	if set, will do basic introspection to forward SSH traffic to this address
```

### Usage with SSH

You can use `forward` to do multiplexing of SSH and HTTP in a quick and dirty way, using very simple protocol introspection. A more robust solution would be [sshttp](https://github.com/stealth/sshttp)

## Usage on Windows

`forward` is wrapped with [go-svc](https://github.com/judwhite/go-svc), enabling it to be run as a Windows service. To add with PowerShell:

```powershell
New-Service -BinaryPathName "C:\path\to\forward.exe -connect 192.168.0.10:80 -listen :80" -Name "port-forward-http"
```
