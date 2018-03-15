# Golang utils

Some simple utility functions, used by most of my projects.

* [MapHelper](map_helper.go): working with maps.
* [Proxy](proxy.go): load/save proxy from/to environment, Gnome, etc.


## Go configuration

```
cat >> ${HOME}/.profile <EOF
export GOROOT=/usr/lib/go-1.8
export PATH=${PATH}:${HOME}/bin:${GOROOT}/bin

for i in ${HOME}/tmp/go ${HOME}/Dropbox/Proyectos/go ; do
	if [ -z "${GOPATH}" ]; then
		export GOPATH=${i}
	else
		export GOPATH=${GOPATH}:${i}
	fi
	export PATH=${PATH}:${i}/bin
done
EOF
```


## Me

Website (in spanish): https://okelet.github.io

Email: okelet@gmail.com
