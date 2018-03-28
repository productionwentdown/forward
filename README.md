
# forward

A simple TCP proxy. Currently used in [AppVenture](https://appventure.nushigh.edu.sg/)'s internal server to port forward from host to a Hyper-V VM. 

## Usage 

```bash
$ ./forward -help
Usage of ./forward:
  -connect string
    	forward to ip and port (default ":8080")
  -listen string
    	listen on ip and port (default ":8081")
```

## Usage on Windows

`forward` is wrapped with [go-svc](https://github.com/judwhite/go-svc), enabling it to be run as a Windows service. To add with PowerShell:

```powershell
New-Service -BinaryPathName "C:\path\to\forward.exe -connect 192.168.0.10:80 -listen :80" -Name "port-forward-http"
```
