package scripts

import "embed"

//go:embed gpu_install.sh cpu_install.sh requirements.txt
var FS embed.FS
