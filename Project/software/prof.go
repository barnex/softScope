package softscope

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

func InitProfiler() {
	if *flag_CPUProf {
		InitCPUProf()
		go func() {
			time.Sleep(1 * time.Minute)
			FlushProf()
		}()
	}
}

func InitCPUProf() {
	// start CPU profile to file
	fname := "softscope.pprof"
	f, err := os.Create(fname)
	check(err)
	err = pprof.StartCPUProfile(f)
	check(err)
	log.Println("writing CPU profile to", fname)

	// at exit: exec go tool pprof to generate SVG output
	AtExit(func() {
		pprof.StopCPUProfile()
		me := procSelfExe()
		outfile := fname + ".svg"
		saveCmdOutput(outfile, "go", "tool", "pprof", "-svg", me, fname)
	})
}

//func InitMemProf() {
//	log.Println("memory profile enabled")
//	AtExit(func() {
//		fname := "softscope.pprof"
//		f, err := os.Create(fname)
//		defer f.Close()
//		check(err)
//		log.Println("writing memory profile to", fname)
//		check(pprof.WriteHeapProfile(f))
//		me := procSelfExe()
//		outfile := fname + ".svg"
//		saveCmdOutput(outfile, "go", "tool", "pprof", "-svg", "--inuse_objects", me, fname)
//	})
//}

// Exec command and write output to outfile.
func saveCmdOutput(outfile string, cmd string, args ...string) {
	log.Println("exec:", cmd, args, ">", outfile)
	out, err := exec.Command(cmd, args...).Output() // TODO: stderr is ignored
	if err != nil {
		log.Printf("exec %v %v: %v: %v", cmd, args, err, string(out))
	}
	check(ioutil.WriteFile(outfile, out, 0666))
}

// path to the executable.
func procSelfExe() string {
	me, err := os.Readlink("/proc/self/exe")
	check(err)
	return me
}

// Functions to be called at program exit
var atexit []func()

// Add a function to be executed at program exit.
func AtExit(cleanup func()) {
	atexit = append(atexit, cleanup)
}

// Runs all functions stacked by AtExit().
func FlushProf() {
	if len(atexit) != 0 {
		log.Println("stopping profiler")
	}
	for _, f := range atexit {
		f()
	}
}
