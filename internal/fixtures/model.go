package fixtures

import (
	"github.com/germainlefebvre4/cvwonder/internal/model"
)

var CvYamlGood01 = []byte(`
person:
  name: John Doe
`)

var CvYamlWithExperienceYears = []byte(`
person:
  name: John Doe
  experience:
    years: 10
`)

var CvYamlWithExperienceSince = []byte(`
person:
  name: Jane Smith
  experience:
    since: 2014
`)

var CvYamlWithExperienceFull = []byte(`
person:
  name: Bob Johnson
  experience:
    years: 10
    since: 2014
`)

var CvYamlWithMissionCompanyLogo = []byte(`
person:
  name: John Doe

career:
  - companyName: Tech Corp
    missions:
      - position: Senior Engineer
        company: Tech Corp
        companyLogo: images/techcorp-logo.webp
        location: Paris, France
        dates: 2020-2024
`)

var CvModelWithMissionCompanyLogo = model.CV{
	Person: model.Person{
		Name: "John Doe",
	},
	Career: []model.Career{
		{
			CompanyName: "Tech Corp",
			Missions: []model.Mission{
				{
					Position:    "Senior Engineer",
					Company:     "Tech Corp",
					CompanyLogo: "images/techcorp-logo.webp",
					Location:    "Paris, France",
					Dates:       "2020-2024",
				},
			},
		},
	},
}

var CvYamlGood02 = []byte(`
---
company:
  name: Zatsit

person:
  name: Germain
  depiction: profile.png
  profession: Platform Engineer
  location: Lille
  citizenship: FR
  email: germain.lefebvre@mycompany.fr
  site: http://germainlefebvre.fr
  phone: +33 6 00 00 00 00

socialNetworks:
  github: germainlefebvre4
  stackoverflow: germainlefebvre4
  linkedin: germainlefebvre4
  twitter: germainlefebvr4

abstract:
  - "I am a Platform Engineer looking for people to share knowledge to each other."
  - "This section can be a multiples lines of text."
  - "This section can be a multiples lines of text again."

career:
  - companyName: Zatsit
    companyLogo: images/zatsit-logo.webp
    duration: 10 mois, aujourd'hui
    missions:
      - position: Platform Engineer
        company: Adeo
        location: Ronchin, France
        dates: 2024, mars - 2024, décembre
        summary: Construire une IDP, plateforme interne de développement, totalement managée pour aider les développeurs à se focaliser sur le code. Sur base du code source, la plateforme provisionne l'infrastructure sous-jacente, les base de données, la construction des artefact et publication sur la registry, le déploiement dans Kubernetes, l'intégration du monitoring avec Datadog et construction des Monitors.
        technologies:
          - ArgoCD
          - Kubernetes
          - K8s Operrator
          - Crossplane
          - Vault
          - Github Actions
          - JFrog Artifactory
          - Backstage
          - Python
          - Golang
        description:
          - Développement de l'operator Kubernetes responsable du provisioning des bases de données
          - Développement des Compositions Crossplane pour provisionner les base de données
          - Développement de l'API de l'IDP en Golang

technicalSkills:
  list: []
  domains:
    - name: Cloud
      competencies:
        - name: AWS
          level: 80
        - name: GCP
          level: 70
        - name: Azure
          level: 40

sideProjects:
  - name: cvwonder
    position: maintainer
    description: A CLI to render your CV from a YAML file.
    link: germainlefebvre4/cvwonder
    type: github
    langs: Go
    color: 3572A5

certifications:
  - companyName: AWS
    certificationName: Solutions Architect Associate
    issuer: Coursera
    date: Mars 2018
    link: https://www.credly.com/badges/dd09dc40-9ef8-43a4-addb-d861d4dadf26/public_url
    badge: images/aws-certified-solutions-architect-associate.png

education:
  - schoolName: IG2I - Centrale
    schoolLogo: images/centrale-lille-logo.webp
    degree: Titre d'ingénieur (BAC+5)
    location: Lens, France
    dates: 2019 - 2014
    link: https://ig2i.centralelille.fr

references:
  - name: Jane Smith
    position: CTO
    company: Tech Innovations
    date: Janvier 2024
    url: https://linkedin.com/in/janesmith
    socialNetworks:
      linkedin: janesmith
      github: janesmith
    description: "Une collaboration exceptionnelle. Germain a su transformer notre infrastructure."
`)

var CvYamlGood03 = []byte(`
another:
  field: value
`)

var CvYamlError01 = []byte(`
person:
  name: John Doe
    depiction: I am a dummy Software Engineer for test.
 wrong: field
`)

var CvModelGood01 = model.CV{
	Person: model.Person{
		Name: "John Doe",
	},
}

var CvModelWithExperienceYears = model.CV{
	Person: model.Person{
		Name: "John Doe",
		Experience: model.Experience{
			Years: 10,
		},
	},
}

var CvModelWithExperienceSince = model.CV{
	Person: model.Person{
		Name: "Jane Smith",
		Experience: model.Experience{
			Since: 2014,
		},
	},
}

var CvModelWithExperienceFull = model.CV{
	Person: model.Person{
		Name: "Bob Johnson",
		Experience: model.Experience{
			Years: 10,
			Since: 2014,
		},
	},
}

var CvModelGood02 = model.CV{
	Company: model.Company{
		Name: "Zatsit",
	},
	Person: model.Person{
		Name:        "Germain",
		Depiction:   "profile.png",
		Profession:  "Platform Engineer",
		Location:    "Lille",
		Citizenship: "FR",
		Email:       "germain.lefebvre@mycompany.fr",
		Site:        "http://germainlefebvre.fr",
		Phone:       "+33 6 00 00 00 00",
	},
	SocialNetworks: model.SocialNetworks{
		Github:        "germainlefebvre4",
		Stackoverflow: "germainlefebvre4",
		Linkedin:      "germainlefebvre4",
		Twitter:       "germainlefebvr4",
	},
	Abstract: []string{
		"I am a Platform Engineer looking for people to share knowledge to each other.",
		"This section can be a multiples lines of text.",
		"This section can be a multiples lines of text again.",
	},
	Career: []model.Career{
		{
			CompanyName: "Zatsit",
			CompanyLogo: "images/zatsit-logo.webp",
			Duration:    "10 mois, aujourd'hui",
			Missions: []model.Mission{
				{
					Position: "Platform Engineer",
					Company:  "Adeo",
					Location: "Ronchin, France",
					Dates:    "2024, mars - 2024, décembre",
					Summary:  "Construire une IDP, plateforme interne de développement, totalement managée pour aider les développeurs à se focaliser sur le code. Sur base du code source, la plateforme provisionne l'infrastructure sous-jacente, les base de données, la construction des artefact et publication sur la registry, le déploiement dans Kubernetes, l'intégration du monitoring avec Datadog et construction des Monitors.",
					Technologies: []string{
						"ArgoCD",
						"Kubernetes",
						"K8s Operrator",
						"Crossplane",
						"Vault",
						"Github Actions",
						"JFrog Artifactory",
						"Backstage",
						"Python",
						"Golang",
					},
					Description: []string{
						"Développement de l'operator Kubernetes responsable du provisioning des bases de données",
						"Développement des Compositions Crossplane pour provisionner les base de données",
						"Développement de l'API de l'IDP en Golang",
					},
				},
			},
		},
	},
	TechnicalSkills: model.TechnicalSkills{
		Domains: []model.Domain{
			{
				Name: "Cloud",
				Competencies: []model.Competency{
					{
						Name:  "AWS",
						Level: 80,
					},
					{
						Name:  "GCP",
						Level: 70,
					},
					{
						Name:  "Azure",
						Level: 40,
					},
				},
			},
		},
	},
	SideProjects: []model.SideProject{
		{
			Name:        "cvwonder",
			Position:    "maintainer",
			Description: "A CLI to render your CV from a YAML file.",
			Link:        "germainlefebvre4/cvwonder",
			Type:        "github",
			Langs:       "Go",
			Color:       "3572A5",
		},
	},
	Certifications: []model.Certification{
		{
			CompanyName:       "AWS",
			CertificationName: "Solutions Architect Associate",
			Issuer:            "Coursera",
			Date:              "Mars 2018",
			Link:              "https://www.credly.com/badges/dd09dc40-9ef8-43a4-addb-d861d4dadf26/public_url",
			Badge:             "images/aws-certified-solutions-architect-associate.png",
		},
	},
	Education: []model.Education{
		{
			SchoolName: "IG2I - Centrale",
			SchoolLogo: "images/centrale-lille-logo.webp",
			Degree:     "Titre d'ingénieur (BAC+5)",
			Location:   "Lens, France",
			Dates:      "2019 - 2014",
			Link:       "https://ig2i.centralelille.fr",
		},
	},
	References: []model.Reference{
		{
			Name:     "Jane Smith",
			Position: "CTO",
			Company:  "Tech Innovations",
			Date:     "Janvier 2024",
			Url:      "https://linkedin.com/in/janesmith",
			SocialNetworks: []model.SocialNetworks{
				{
					Linkedin: "janesmith",
					Github:   "janesmith",
				},
			},
			Description: "Une collaboration exceptionnelle. Germain a su transformer notre infrastructure.",
		},
	},
}

var CvModelGood03 = model.CV{}

var CvModelError01 = model.CV{}

var CvHtmlTemplate01 = []byte(`{{ .Person.Name }}`)

var CvHtmlGood01 = []byte(`John Doe`)

var CvHtmlTemplateFunctionsTemplateInc01 = []byte(`{{ inc 1 }}`)
var CvHtmlTemplateFunctionsGeneratedInc01 = []byte(`2`)

var CvHtmlTemplateFunctionsTemplateDec01 = []byte(`{{ dec 2 }}`)
var CvHtmlTemplateFunctionsGeneratedDec01 = []byte(`1`)

var CvHtmlTemplateFunctionsTemplateList01 = []byte(`{{ list "a" "b" "c" }}`)
var CvHtmlTemplateFunctionsGeneratedList01 = []byte(`[a b c]`)

var CvHtmlTemplateFunctionsTemplateJoin01 = []byte(`{{ join (list "a" "b" "c") " " }}`)
var CvHtmlTemplateFunctionsGeneratedJoin01 = []byte(`a b c`)

var CvHtmlTemplateFunctionsTemplateSplit01 = []byte(`{{ split "a b c" " " }}`)
var CvHtmlTemplateFunctionsGeneratedSplit01 = []byte(`[a b c]`)

var CvHtmlTemplateFunctionsTemplateTrim01 = []byte(`{{ trim "  a b c  " }}`)
var CvHtmlTemplateFunctionsGeneratedTrim01 = []byte(`a b c`)

var CvHtmlTemplateFunctionsTemplateLower01 = []byte(`{{ lower "A B C" }}`)
var CvHtmlTemplateFunctionsGeneratedLower01 = []byte(`a b c`)

var CvHtmlTemplateFunctionsTemplateUpper01 = []byte(`{{ upper "a b c" }}`)
var CvHtmlTemplateFunctionsGeneratedUpper01 = []byte(`A B C`)

var CvHtmlTemplateFunctionsTemplateReplace01 = []byte(`{{ replace "a b c" " " "-" }}`)
var CvHtmlTemplateFunctionsGeneratedReplace01 = []byte(`a-b-c`)
