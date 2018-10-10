package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"pflag"
)

/* selpg_args */
type selpg_args struct {
	start_page int
	end_page int
	in_filename string
	page_len int
	page_type bool
	print_dest string
}

/* check for visited arg */
var visit map[string]bool = make(map[string]bool)

func main(){
	args := new(selpg_args)
	getArgs(args)
	checkAgrs(args)
	inputbuf := getInput(args)
	output(inputbuf, args)
}

func getArgs(args *selpg_args){
	/* Init Arg_s */
	pflag.IntVarP(&args.start_page, "start","s", -1, "Start Page")
	pflag.IntVarP(&args.end_page, "end", "e", -1, "End Page")
	pflag.BoolVarP(&args.page_type, "f", "f",false, "Find page breaks")
	pflag.IntVarP(&args.page_len, "leap", "l", 72, "The count of rows of each page")
	pflag.StringVarP(&args.print_dest, "destination", "d", "", "Print Destination")
	pflag.Parse()
	if len(pflag.Args()) == 1 {
		args.in_filename = pflag.Arg(0)
	}
	pflag.Visit(func(f *pflag.Flag) {
		visit[f.Name] = true
	})
}

func checkAgrs(args *selpg_args) {
	switch {
	case !visit["start"] :
		fmt.Print("-s --start miss\n")
	case !visit["end"] :
		fmt.Print("-e --end miss\n")
	case args.start_page <= 0:
		fmt.Print("-s --start Invalid Page Number\n")
	case args.end_page <= 0:
		fmt.Print("-e --end Invalid Page Number\n")
	case args.end_page < args.start_page:
		fmt.Print("Start Page must smaller than End Page\n")
	case visit["f"] && visit["leap"]:
		fmt.Print("can't use two mode at the same time\n")
	}
}

func getInput(args *selpg_args) *bufio.Reader{
	if len(args.in_filename) == 0{
		/* from stdin */
		return bufio.NewReader(os.Stdin)
	} else{
		/* from input_file */
		files, err := os.OpenFile(args.in_filename, 0, 0)
		controlErr(err)
		return bufio.NewReader(files)
	}
}

func controlErr(err error) {
	if err != nil {
		panic(err)
	}
}

func output(inputbuf *bufio.Reader, args *selpg_args) {
	var out string
	if(args.page_type) {
		/* -f */
		page := 0
		for {
			line, err := inputbuf.ReadString('\f')
			if err != nil && err != io.EOF {
				controlErr(err)
			}
			page++
			if args.start_page <= page && page <= args.end_page{
				out += line
			}
			if err == io.EOF {
				break
			}
		}

	}else {
		/* -l */
		lineCount := 0
		page := 1
		for {
			line, err := inputbuf.ReadString('\n')
			if err != nil && err != io.EOF {
				controlErr(err)
			}
			lineCount++
			if args.start_page <= page && page <= args.end_page{
				out += line
			}
			if lineCount >= args.page_len{
				page++
				lineCount = 0
			}
			if err == io.EOF {
				break
			}
		}
	}
	if args.print_dest == "" {
		os.Stdout.WriteString(out)
	}else {
		cmd := exec.Command("lp", "-d" + args.print_dest)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			controlErr(err)
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, out)
		}()
	}

}