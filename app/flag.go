package app

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

// NamedFlagSets stores named flag sets in the order of calling FlagSet
type NamedFlagSets struct {
	Order    []string
	FlagSets map[string]*pflag.FlagSet
}

// FlagSet returns the FlagSets by the given name
// and adds the flag to the ordered name list if it is not in there
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}

	if _, ok := nfs.FlagSets[name]; !ok {
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}

	return nfs.FlagSets[name]
}

func InitFlags(flags *pflag.FlagSet) {
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	// if flag package has been used
	// flags.AddGoFlagSet()
}

// WordSepNormalizeFunc changes all flags that contain "_" separators.
func WordSepNormalizeFunc(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}

	return pflag.NormalizedName(name)
}

// AddFlags
func AddFlags(flagName string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(flagName))
}

// PrintSections prints the given names flagsets info in sections
func PrintSections(w io.Writer, fs NamedFlagSets, cols int) {
	for _, name := range fs.Order {
		fs := fs.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFs := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFs.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFs.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFs.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, buf.String())
		}
	}
}

// printFalgs prints flag in flagsets
func printFalgs(fs *pflag.FlagSet) {
	fs.VisitAll(func(f *pflag.Flag) {
		fmt.Printf("FLAG: -%s=%q", f.Name, f.Value)
	})
}
