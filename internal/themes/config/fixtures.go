package theme_config

var ThemeConfigYamlGood01 = []byte(`
name: Test01
slug: test01
`)

var ThemeConfigYamlGood02 = []byte(`
name: Test02
slug: test02
description: Test02
author: Test02
`)

var ThemeConfigYamlGood03 = []byte(`
name: Test03
slug: test03
description: Test03
author: Test03
minimumVersion: 0.3.1
`)

var ThemeConfigModelGood01 = ThemeConfig{
	Name: "Test01",
	Slug: "test01",
}

var ThemeConfigModelGood02 = ThemeConfig{
	Name:        "Test02",
	Slug:        "test02",
	Description: "Test02",
	Author:      "Test02",
}

var ThemeConfigModelGood03 = ThemeConfig{
	Name:           "Test03",
	Slug:           "test03",
	Description:    "Test03",
	Author:         "Test03",
	MinimumVersion: "0.3.1",
}
