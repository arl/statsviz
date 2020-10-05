// memStatsDoc holds documentation for MemStats fields
var memStatsDoc = (function () {
  const docs = {
    "Alloc": "Alloc is bytes of allocated heap objects.\n\nThis is the same as HeapAlloc (see below).\n",
    "BuckHashSys": "BuckHashSys is bytes of memory in profiling bucket hash tables.\n",
    "BySize": "BySize reports per-size class allocation statistics.\n\nBySize[N] gives statistics for allocations of size S where\nBySize[N-1].Size \u003c S ≤ BySize[N].Size.\n\nThis does not report allocations larger than BySize[60].Size.\n",
    "DebugGC": "DebugGC is currently unused.\n",
    "EnableGC": "EnableGC indicates that GC is enabled. It is always true,\neven if GOGC=off.\n",
    "Frees": "Frees is the cumulative count of heap objects freed.\n",
    "GCCPUFraction": "GCCPUFraction is the fraction of this program's available\nCPU time used by the GC since the program started.\n\nGCCPUFraction is expressed as a number between 0 and 1,\nwhere 0 means GC has consumed none of this program's CPU. A\nprogram's available CPU time is defined as the integral of\nGOMAXPROCS since the program started. That is, if\nGOMAXPROCS is 2 and a program has been running for 10\nseconds, its \"available CPU\" is 20 seconds. GCCPUFraction\ndoes not include CPU time used for write barrier activity.\n\nThis is the same as the fraction of CPU reported by\nGODEBUG=gctrace=1.\n",
    "GCSys": "GCSys is bytes of memory in garbage collection metadata.\n",
    "HeapAlloc": "HeapAlloc is bytes of allocated heap objects.\n\n\"Allocated\" heap objects include all reachable objects, as\nwell as unreachable objects that the garbage collector has\nnot yet freed. Specifically, HeapAlloc increases as heap\nobjects are allocated and decreases as the heap is swept\nand unreachable objects are freed. Sweeping occurs\nincrementally between GC cycles, so these two processes\noccur simultaneously, and as a result HeapAlloc tends to\nchange smoothly (in contrast with the sawtooth that is\ntypical of stop-the-world garbage collectors).\n",
    "HeapIdle": "HeapIdle is bytes in idle (unused) spans.\n\nIdle spans have no objects in them. These spans could be\n(and may already have been) returned to the OS, or they can\nbe reused for heap allocations, or they can be reused as\nstack memory.\n\nHeapIdle minus HeapReleased estimates the amount of memory\nthat could be returned to the OS, but is being retained by\nthe runtime so it can grow the heap without requesting more\nmemory from the OS. If this difference is significantly\nlarger than the heap size, it indicates there was a recent\ntransient spike in live heap size.\n",
    "HeapInuse": "HeapInuse is bytes in in-use spans.\n\nIn-use spans have at least one object in them. These spans\ncan only be used for other objects of roughly the same\nsize.\n\nHeapInuse minus HeapAlloc estimates the amount of memory\nthat has been dedicated to particular size classes, but is\nnot currently being used. This is an upper bound on\nfragmentation, but in general this memory can be reused\nefficiently.\n",
    "HeapObjects": "HeapObjects is the number of allocated heap objects.\n\nLike HeapAlloc, this increases as objects are allocated and\ndecreases as the heap is swept and unreachable objects are\nfreed.\n",
    "HeapReleased": "HeapReleased is bytes of physical memory returned to the OS.\n\nThis counts heap memory from idle spans that was returned\nto the OS and has not yet been reacquired for the heap.\n",
    "HeapSys": "HeapSys is bytes of heap memory obtained from the OS.\n\nHeapSys measures the amount of virtual address space\nreserved for the heap. This includes virtual address space\nthat has been reserved but not yet used, which consumes no\nphysical memory, but tends to be small, as well as virtual\naddress space for which the physical memory has been\nreturned to the OS after it became unused (see HeapReleased\nfor a measure of the latter).\n\nHeapSys estimates the largest size the heap has had.\n",
    "LastGC": "LastGC is the time the last garbage collection finished, as\nnanoseconds since 1970 (the UNIX epoch).\n",
    "Lookups": "Lookups is the number of pointer lookups performed by the\nruntime.\n\nThis is primarily useful for debugging runtime internals.\n",
    "MCacheInuse": "MCacheInuse is bytes of allocated mcache structures.\n",
    "MCacheSys": "MCacheSys is bytes of memory obtained from the OS for\nmcache structures.\n",
    "MSpanInuse": "MSpanInuse is bytes of allocated mspan structures.\n",
    "MSpanSys": "MSpanSys is bytes of memory obtained from the OS for mspan\nstructures.\n",
    "Mallocs": "Mallocs is the cumulative count of heap objects allocated.\nThe number of live objects is Mallocs - Frees.\n",
    "NextGC": "NextGC is the target heap size of the next GC cycle.\n\nThe garbage collector's goal is to keep HeapAlloc ≤ NextGC.\nAt the end of each GC cycle, the target for the next cycle\nis computed based on the amount of reachable data and the\nvalue of GOGC.\n",
    "NumForcedGC": "NumForcedGC is the number of GC cycles that were forced by\nthe application calling the GC function.\n",
    "NumGC": "NumGC is the number of completed GC cycles.\n",
    "OtherSys": "OtherSys is bytes of memory in miscellaneous off-heap\nruntime allocations.\n",
    "PauseEnd": "PauseEnd is a circular buffer of recent GC pause end times,\nas nanoseconds since 1970 (the UNIX epoch).\n\nThis buffer is filled the same way as PauseNs. There may be\nmultiple pauses per GC cycle; this records the end of the\nlast pause in a cycle.\n",
    "PauseNs": "PauseNs is a circular buffer of recent GC stop-the-world\npause times in nanoseconds.\n\nThe most recent pause is at PauseNs[(NumGC+255)%256]. In\ngeneral, PauseNs[N%256] records the time paused in the most\nrecent N%256th GC cycle. There may be multiple pauses per\nGC cycle; this is the sum of all pauses during a cycle.\n",
    "PauseTotalNs": "PauseTotalNs is the cumulative nanoseconds in GC\nstop-the-world pauses since the program started.\n\nDuring a stop-the-world pause, all goroutines are paused\nand only the garbage collector can run.\n",
    "StackInuse": "StackInuse is bytes in stack spans.\n\nIn-use stack spans have at least one stack in them. These\nspans can only be used for other stacks of the same size.\n\nThere is no StackIdle because unused stack spans are\nreturned to the heap (and hence counted toward HeapIdle).\n",
    "StackSys": "StackSys is bytes of stack memory obtained from the OS.\n\nStackSys is StackInuse, plus any memory obtained directly\nfrom the OS for OS thread stacks (which should be minimal).\n",
    "Sys": "Sys is the total bytes of memory obtained from the OS.\n\nSys is the sum of the XSys fields below. Sys measures the\nvirtual address space reserved by the Go runtime for the\nheap, stacks, and other internal data structures. It's\nlikely that not all of the virtual address space is backed\nby physical memory at any given moment, though in general\nit all was at some point.\n",
    "TotalAlloc": "TotalAlloc is cumulative bytes allocated for heap objects.\n\nTotalAlloc increases as heap objects are allocated, but\nunlike Alloc and HeapAlloc, it does not decrease when\nobjects are freed.\n"
  };
  return docs;
}());
