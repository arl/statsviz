package plot

var (
	idxcgogotocalls                 int
	idxcpuclassesgcmarkassist       int
	idxcpuclassesgcmarkdedicated    int
	idxcpuclassesgcmarkidle         int
	idxcpuclassesgcpause            int
	idxcpuclassesgctotal            int
	idxcpuclassesidle               int
	idxcpuclassesscavengeassist     int
	idxcpuclassesscavengebackground int
	idxcpuclassesscavengtetotal     int
	idxcpuclassestotal              int
	idxcpuclassesuser               int
)

func init() {
	idxcgogotocalls = mustidx("/cgo/go-to-c-calls:calls")
	idxcpuclassesgcmarkassist = mustidx("/cpu/classes/gc/mark/assist:cpu-seconds")
	idxcpuclassesgcmarkdedicated = mustidx("/cpu/classes/gc/mark/dedicated:cpu-seconds")
	idxcpuclassesgcmarkidle = mustidx("/cpu/classes/gc/mark/idle:cpu-seconds")
	idxcpuclassesgcpause = mustidx("/cpu/classes/gc/pause:cpu-seconds")
	idxcpuclassesgctotal = mustidx("/cpu/classes/gc/total:cpu-seconds")
	idxcpuclassesidle = mustidx("/cpu/classes/idle:cpu-seconds")
	idxcpuclassesscavengeassist = mustidx("/cpu/classes/scavenge/assist:cpu-seconds")
	idxcpuclassesscavengebackground = mustidx("/cpu/classes/scavenge/background:cpu-seconds")
	idxcpuclassesscavengtetotal = mustidx("/cpu/classes/scavenge/total:cpu-seconds")
	idxcpuclassestotal = mustidx("/cpu/classes/total:cpu-seconds")
	idxcpuclassesuser = mustidx("/cpu/classes/user:cpu-seconds")
}
