package gs2parser

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCPPParityBaselines(t *testing.T) {
	base := os.Getenv("GS2_TESTS_DIR")
	if base == "" {
		if cpp := os.Getenv("GS2_CPP_REPO"); cpp != "" {
			base = filepath.Join(cpp, "tests")
		} else {
			base = "tests"
		}
	}
	var cases []string
	if os.Getenv("GS2_ALL_BASELINES") == "1" {
		err := filepath.WalkDir(filepath.Join(base, "baselines"), func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".bytecode" {
				return err
			}
			rel, err := filepath.Rel(filepath.Join(base, "baselines"), path)
			if err != nil {
				return err
			}
			cases = append(cases, rel[:len(rel)-len(".bytecode")])
			return nil
		})
		if err != nil {
			t.Skip(err)
		}
	} else {
		cases = []string{
			"basic/01_variables",
			"basic/02_constants",
			"basic/03_data_types",
			"expressions/01_arithmetic",
			"expressions/02_comparison",
			"expressions/03_logical",
			"expressions/04_bitwise",
			"expressions/05_assignment",
			"statements/01_conditionals",
			"statements/02_loops",
			"statements/03_switch",
			"statements/04_with",
			"functions/01_basic_functions",
			"functions/02_recursion",
			"functions/03_lambdas",
			"classes/01_objects",
			"classes/02_arrays",
			"advanced/g2k1/weaponTestGS2",
			"advanced/loginserver/weapon-Staff_GraalShop",
		}
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			src, err := os.ReadFile(filepath.Join(base, "scripts", name+".gs2"))
			if err != nil {
				t.Skip(err)
			}
			want, err := os.ReadFile(filepath.Join(base, "baselines", name+".bytecode"))
			if err != nil {
				t.Skip(err)
			}
			got, err := Compile(string(src))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got.Bytecode, want) {
				t.Fatalf("bytecode mismatch: got %d bytes, want %d bytes", len(got.Bytecode), len(want))
			}
		})
	}
}
