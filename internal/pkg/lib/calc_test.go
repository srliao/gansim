package lib

import (
	"io/ioutil"
	"log"
	"math"
	"testing"

	"gopkg.in/yaml.v2"
)

type testCfg struct {
	Profile   Profile           `yaml:"Profile"`
	Artifacts map[Slot]Artifact `yaml:"Artifacts"`
}

func TestCalc(t *testing.T) {

	var source []byte
	var cfg testCfg
	var err error

	source, err = ioutil.ReadFile("./test/ganyu_lvl88fatui.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	var r []DmgResult

	r = Calc(cfg.Profile, cfg.Artifacts, false)

	checkt(r[0].Normal, 3668, 0.001, t)
	checkt(r[0].Crit, 10674, 0.001, t)
	log.Println(r)

	//should be 3668 for lvl 88??

	source, err = ioutil.ReadFile("./test/ganyu_melt.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	r = Calc(cfg.Profile, cfg.Artifacts, true)

	checkt(r[0].Normal, 2934, 0.001, t)
	checkt(r[0].Crit, 8540, 0.001, t)
	checkt(r[1].Normal, 5682, 0.001, t)
	checkt(r[1].Crit, 16535, 0.001, t)
	log.Println(r)

	//diff should be under 1%

}

func checkt(got, expected, tolerance float64, t *testing.T) {
	if math.Abs((got-expected)/expected) > tolerance {
		t.Errorf("expected %v got %.4f, tol % .4f", expected, got, math.Abs((got-expected)/expected))
	}
}
