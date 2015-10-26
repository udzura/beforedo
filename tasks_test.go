package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	conf, _ := ioutil.TempFile("", "beforedo--testing--")
	defer os.Remove(conf.Name())
	conf.Write([]byte(`
- task: mysql.server start
  port: 3306
- task: memcached -p 11211 -m 64m -d
  port: 11211
- task: bundle install --path vendor/bundle
  success: bundle check
- task: cp config/database.yml.sample config/database.yml
  file: config/database.yml
- task: bundle exec rake db:migrate
  always: true
- task: bundle exec rails s
  front: true
`))
	conf.Close()

	tasks, _ := ParseConfig(conf.Name())

	if l := len(tasks); l != 6 {
		t.Fatalf("Config should have 6 tasks, got %d", l)
	}

	if t1 := tasks[0]; t1.Command != "mysql.server start" || t1.Port != 3306 {
		t.Fatalf("Parser failed: %+v", t1)
	}

	if t2 := tasks[1]; t2.Command != "memcached -p 11211 -m 64m -d" || t2.Port != 11211 {
		t.Fatalf("Parser failed: %+v", t2)
	}

	if t3 := tasks[2]; t3.Command != "bundle install --path vendor/bundle" || t3.SuccessCommand != "bundle check" {
		t.Fatalf("Parser failed: %+v", t3)
	}

	if t4 := tasks[3]; t4.Command != "cp config/database.yml.sample config/database.yml" || t4.File != "config/database.yml" {
		t.Fatalf("Parser failed: %+v", t4)
	}

	if t5 := tasks[4]; t5.Command != "bundle exec rake db:migrate" || !t5.Always {
		t.Fatalf("Parser failed: %+v", t5)
	}

	if t6 := tasks[5]; t6.Command != "bundle exec rails s" || !t6.Front {
		t.Fatalf("Parser failed: %+v", t6)
	}
}
