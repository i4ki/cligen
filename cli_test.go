package cligen

import (
	"fmt"
	"testing"

	"github.com/madlambda/spells/assert"
)

type parseTestcase struct {
	in   string
	want []Cli
}

func TestParse(t *testing.T) {
	testcases := []parseTestcase{
		{
			in: `package s
// Copy files and directories.
func Copy(src string, dst string, recursive bool) error {return nil}`,
			want: []Cli{
				{
					Name: "Copy",
					Desc: "Copy files and directories.",
					Args: []Arg{
						{
							Type: "string",
							Name: "src",
							Desc: "",
						},
						{
							Type: "string",
							Name: "dst",
							Desc: "",
						},
					},
					Flags: []Flag{
						{
							Name: "recursive",
							Desc: "",
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		clis, err := Parse("test.go", tc.in)
		assert.NoError(t, err, "failed to parse")

		assert.EqualInts(t, len(tc.want), len(clis), "clis len")

		for i := 0; i < len(tc.want); i++ {
			cli := clis[i]
			want := tc.want[i]

			help, err := cli.Help()
			assert.NoError(t, err, "cli help")

			fmt.Printf("'%s'\n", help)

			assert.EqualStrings(t, want.Name, cli.Name, "cli name")
			assert.EqualStrings(t, want.Desc, cli.Desc, "cli desc")

			assert.EqualInts(t, len(want.Flags), len(cli.Flags), "flags len")

			for j := 0; j < len(want.Flags); j++ {
				wantFlag := want.Flags[j]
				gotFlag := cli.Flags[j]
				assert.EqualStrings(t, wantFlag.Name, gotFlag.Name, "flag name")
				assert.EqualStrings(t, wantFlag.Desc, gotFlag.Desc, "flag desc")
			}

			assert.EqualInts(t, len(want.Args), len(cli.Args), "args len")

			for j := 0; j < len(want.Args); j++ {
				wantArg := want.Args[j]
				gotArg := cli.Args[j]
				assert.EqualStrings(t, wantArg.Name, gotArg.Name, "arg name")
				assert.EqualStrings(t, wantArg.Type, gotArg.Type, "arg type")
				assert.EqualStrings(t, wantArg.Desc, gotArg.Desc, "arg desc")
			}
		}
	}
}
