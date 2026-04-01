package cvinit

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/huh/v2"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/goccy/go-yaml"
)

// writePartial marshals cv and writes it to path, creating or overwriting the file.
func writePartial(cv model.CV, path string) error {
	data, err := yaml.Marshal(cv)
	if err != nil {
		return fmt.Errorf("marshal cv: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write cv: %w", err)
	}
	return nil
}

// RunWizard runs the interactive CV wizard and writes the result to outputFile.
// Returns an error if outputFile already exists.
func RunWizard(outputFile string) error {
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("file already exists: %s (remove it or use --output-file to choose a different name)", outputFile)
	}

	var cv model.CV

	// ── Step 1-2: Welcome + output-file confirmation ─────────────────────────
	if err := runWelcome(&outputFile); err != nil {
		return err
	}

	// ── Step 3: Company ───────────────────────────────────────────────────────
	if err := runCompany(&cv); err != nil {
		return err
	}
	_ = writePartial(cv, outputFile)

	// ── Step 4: Person (mandatory) ────────────────────────────────────────────
	if err := runPerson(&cv); err != nil {
		return err
	}
	_ = writePartial(cv, outputFile)

	// ── Step 5: Social Networks ───────────────────────────────────────────────
	if err := runSocialNetworks(&cv); err != nil {
		return err
	}
	_ = writePartial(cv, outputFile)

	// ── Step 6: Abstract ──────────────────────────────────────────────────────
	if err := runAbstract(&cv); err != nil {
		return err
	}
	_ = writePartial(cv, outputFile)

	// ── Step 7: Career ────────────────────────────────────────────────────────
	if err := runCareer(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 8: Technical Skills ──────────────────────────────────────────────
	if err := runTechnicalSkills(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 9: Education ─────────────────────────────────────────────────────
	if err := runEducation(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 10: Certifications ───────────────────────────────────────────────
	if err := runCertifications(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 11: Languages ────────────────────────────────────────────────────
	if err := runLanguages(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 12: Side Projects ────────────────────────────────────────────────
	if err := runSideProjects(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 13: References ───────────────────────────────────────────────────
	if err := runReferences(&cv, outputFile); err != nil {
		return err
	}

	// ── Step 14: Final confirmation + write ───────────────────────────────────
	return runFinalize(cv, outputFile)
}

// ── Section runners ───────────────────────────────────────────────────────────

func runWelcome(outputFile *string) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Welcome to CV Wonder!").
				Description("This wizard will help you create your cv.yml file step by step.\nYou can skip optional sections. Press Ctrl+C at any time — your progress is saved after each section."),
			huh.NewInput().
				Title("Output filename").
				Description("Where to save your CV YAML file.").
				Value(outputFile),
		),
	).Run()
}

func runCompany(cv *model.CV) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Company").Description("Optional: fill if generating a CV for a specific company."),
			huh.NewInput().
				Title("Company name").
				Value(&cv.Company.Name),
			huh.NewInput().
				Title("Company logo").
				Description("Relative path to the logo image (e.g. images/logo.png). Leave blank to fill later.").
				Value(&cv.Company.Logo),
		),
	).Run()
}

func runPerson(cv *model.CV) error {
	yearsStr := ""
	sinceStr := ""

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Personal Information").Description("Required: your core details."),
			huh.NewInput().
				Title("Full name").
				Value(&cv.Person.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Profession / Job title").
				Value(&cv.Person.Profession).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("profession is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Profile picture").
				Description("Relative path to your photo (e.g. images/photo.jpg). Leave blank to fill later.").
				Value(&cv.Person.Depiction),
			huh.NewInput().
				Title("Location").
				Description("e.g. Paris, France").
				Value(&cv.Person.Location),
			huh.NewInput().
				Title("Citizenship").
				Description("e.g. FR").
				Value(&cv.Person.Citizenship),
			huh.NewInput().
				Title("Email").
				Value(&cv.Person.Email),
			huh.NewInput().
				Title("Phone").
				Value(&cv.Person.Phone),
			huh.NewInput().
				Title("Website").
				Value(&cv.Person.Site),
		),
		huh.NewGroup(
			huh.NewNote().Title("Experience").Description("Optional: your years of experience."),
			huh.NewInput().
				Title("Years of experience").
				Description("Leave blank to skip.").
				Value(&yearsStr).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					n, err := strconv.Atoi(strings.TrimSpace(s))
					if err != nil || n < 0 {
						return fmt.Errorf("must be a positive number")
					}
					return nil
				}),
			huh.NewInput().
				Title("Working since (year)").
				Description("e.g. 2015 — leave blank to skip.").
				Value(&sinceStr).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					n, err := strconv.Atoi(strings.TrimSpace(s))
					if err != nil || n < 1900 {
						return fmt.Errorf("must be a valid year (e.g. 2015)")
					}
					return nil
				}),
		),
	).Run()
	if err != nil {
		return err
	}

	if yearsStr != "" {
		cv.Person.Experience.Years, _ = strconv.Atoi(strings.TrimSpace(yearsStr))
	}
	if sinceStr != "" {
		cv.Person.Experience.Since, _ = strconv.Atoi(strings.TrimSpace(sinceStr))
	}
	return nil
}

func runSocialNetworks(cv *model.CV) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add social networks?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("GitHub username").Value(&cv.SocialNetworks.Github),
			huh.NewInput().Title("LinkedIn username").Value(&cv.SocialNetworks.Linkedin),
			huh.NewInput().Title("Stack Overflow username").Value(&cv.SocialNetworks.Stackoverflow),
			huh.NewInput().Title("Twitter / X username").Value(&cv.SocialNetworks.Twitter),
			huh.NewInput().Title("Bluesky handle").Value(&cv.SocialNetworks.Bluesky),
		),
	).Run()
}

func runAbstract(cv *model.CV) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add a professional summary?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	var raw string
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Professional summary").
				Description("Enter each paragraph on a new line.").
				Value(&raw),
		),
	).Run(); err != nil {
		return err
	}
	cv.Abstract = splitLines(raw)
	return nil
}

func runCareer(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add career experience?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		career, err := collectCareerEntry()
		if err != nil {
			return err
		}
		cv.Career = append(cv.Career, career)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another company?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectCareerEntry() (model.Career, error) {
	var career model.Career

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Career Entry").Description("Company details"),
			huh.NewInput().
				Title("Company name").
				Value(&career.CompanyName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("company name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Company logo").
				Description("Relative path to logo image. Leave blank to fill later.").
				Value(&career.CompanyLogo),
			huh.NewInput().
				Title("Duration").
				Description("e.g. 2 years, 6 months").
				Value(&career.Duration),
		),
	).Run(); err != nil {
		return career, err
	}

	for {
		mission, err := collectMission()
		if err != nil {
			return career, err
		}
		career.Missions = append(career.Missions, mission)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another mission at this company?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return career, err
		}
	}
}

func collectMission() (model.Mission, error) {
	var mission model.Mission
	var techRaw, descRaw string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Mission").Description("Position details"),
			huh.NewInput().
				Title("Position / job title").
				Value(&mission.Position).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("position is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Client / department").
				Description("Leave blank if same as company.").
				Value(&mission.Company),
			huh.NewInput().
				Title("Client logo").
				Description("Relative path to client logo. Leave blank to fill later.").
				Value(&mission.CompanyLogo),
			huh.NewInput().
				Title("Location").
				Value(&mission.Location),
			huh.NewInput().
				Title("Dates").
				Description("e.g. 2022 - 2024").
				Value(&mission.Dates),
		),
		huh.NewGroup(
			huh.NewText().
				Title("Mission summary").
				Description("A short paragraph describing the role and context.").
				Value(&mission.Summary),
			huh.NewInput().
				Title("Technologies").
				Description("Comma-separated list, e.g. Go, Docker, Kubernetes").
				Value(&techRaw),
			huh.NewText().
				Title("Description").
				Description("Key achievements and responsibilities — one per line.").
				Value(&descRaw),
			huh.NewInput().
				Title("Project name").
				Description("Optional project name.").
				Value(&mission.Project),
		),
	).Run()
	if err != nil {
		return mission, err
	}

	mission.Technologies = splitComma(techRaw)
	mission.Description = splitLines(descRaw)
	return mission, nil
}

func runTechnicalSkills(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add technical skills?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		domain, err := collectDomain()
		if err != nil {
			return err
		}
		cv.TechnicalSkills.Domains = append(cv.TechnicalSkills.Domains, domain)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another skill domain?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectDomain() (model.Domain, error) {
	var domain model.Domain

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Domain name").
				Description("e.g. Backend Development, Infrastructure").
				Value(&domain.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("domain name is required")
					}
					return nil
				}),
		),
	).Run(); err != nil {
		return domain, err
	}

	for {
		comp, err := collectCompetency()
		if err != nil {
			return domain, err
		}
		domain.Competencies = append(domain.Competencies, comp)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another skill in this domain?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return domain, err
		}
	}
}

func collectCompetency() (model.Competency, error) {
	var comp model.Competency
	var levelStr string

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Skill name").
				Value(&comp.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("skill name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Level").
				Description("1 (beginner) to 5 (expert)").
				Value(&levelStr).
				Validate(validateCompetencyLevel),
		),
	).Run()
	if err != nil {
		return comp, err
	}

	comp.Level, _ = strconv.Atoi(strings.TrimSpace(levelStr))
	return comp, nil
}

func runEducation(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add education?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		edu, err := collectEducation()
		if err != nil {
			return err
		}
		cv.Education = append(cv.Education, edu)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another education entry?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectEducation() (model.Education, error) {
	var edu model.Education

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("School / University name").
				Value(&edu.SchoolName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("school name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("School logo").
				Description("Relative path to logo image. Leave blank to fill later.").
				Value(&edu.SchoolLogo),
			huh.NewInput().
				Title("Degree").
				Description("e.g. Bachelor of Science in Computer Science").
				Value(&edu.Degree),
			huh.NewInput().
				Title("Location").
				Value(&edu.Location),
			huh.NewInput().
				Title("Dates").
				Description("e.g. 2015 - 2019").
				Value(&edu.Dates),
			huh.NewInput().
				Title("Link").
				Description("School or programme URL (optional).").
				Value(&edu.Link),
		),
	).Run()

	return edu, err
}

func runCertifications(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add certifications?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		cert, err := collectCertification()
		if err != nil {
			return err
		}
		cv.Certifications = append(cv.Certifications, cert)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another certification?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectCertification() (model.Certification, error) {
	var cert model.Certification

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Certification name").
				Value(&cert.CertificationName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("certification name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Issuing company / organisation").
				Value(&cert.CompanyName),
			huh.NewInput().
				Title("Issuer").
				Value(&cert.Issuer),
			huh.NewInput().
				Title("Date").
				Description("e.g. 2023").
				Value(&cert.Date),
			huh.NewInput().
				Title("Link").
				Description("Certification URL (optional).").
				Value(&cert.Link),
			huh.NewInput().
				Title("Badge image").
				Description("Relative path to badge image. Leave blank to fill later.").
				Value(&cert.Badge),
		),
	).Run()

	return cert, err
}

func runLanguages(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add languages?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		lang, err := collectLanguage()
		if err != nil {
			return err
		}
		cv.Languages = append(cv.Languages, lang)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another language?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectLanguage() (model.Language, error) {
	var lang model.Language

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Language").
				Value(&lang.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("language name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Level").
				Description("e.g. Native, Fluent, Intermediate").
				Value(&lang.Level),
		),
	).Run()

	return lang, err
}

func runSideProjects(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add side projects?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		proj, err := collectSideProject()
		if err != nil {
			return err
		}
		cv.SideProjects = append(cv.SideProjects, proj)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another side project?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectSideProject() (model.SideProject, error) {
	var proj model.SideProject

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Value(&proj.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("project name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Your role / position").
				Value(&proj.Position),
			huh.NewInput().
				Title("Description").
				Value(&proj.Description),
			huh.NewInput().
				Title("Link").
				Description("URL to the project (optional).").
				Value(&proj.Link),
			huh.NewInput().
				Title("Type").
				Description("e.g. open-source, personal").
				Value(&proj.Type),
			huh.NewInput().
				Title("Languages / technologies").
				Description("e.g. Go, TypeScript").
				Value(&proj.Langs),
			huh.NewInput().
				Title("Color").
				Description("Optional hex color, e.g. #3572A5").
				Value(&proj.Color),
		),
	).Run()

	return proj, err
}

func runReferences(cv *model.CV, outputFile string) error {
	var add bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add references?").
				Value(&add),
		),
	).Run(); err != nil || !add {
		return err
	}

	for {
		ref, err := collectReference()
		if err != nil {
			return err
		}
		cv.References = append(cv.References, ref)
		_ = writePartial(*cv, outputFile)

		var another bool
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another reference?").
					Value(&another),
			),
		).Run(); err != nil || !another {
			return err
		}
	}
}

func collectReference() (model.Reference, error) {
	var ref model.Reference

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&ref.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Position").
				Value(&ref.Position),
			huh.NewInput().
				Title("Company").
				Value(&ref.Company),
			huh.NewInput().
				Title("Date").
				Value(&ref.Date),
			huh.NewInput().
				Title("URL").
				Description("LinkedIn or website URL (optional).").
				Value(&ref.Url),
			huh.NewText().
				Title("Description").
				Description("Optional quote or description.").
				Value(&ref.Description),
		),
	).Run()

	return ref, err
}

func runFinalize(cv model.CV, outputFile string) error {
	var confirm bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Generate %s?", outputFile)).
				Description("This will write your CV YAML file to disk.").
				Affirmative("Generate").
				Negative("Cancel").
				Value(&confirm),
		),
	).Run(); err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Cancelled. No file was written.")
		return nil
	}

	if err := writePartial(cv, outputFile); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputFile, err)
	}

	fmt.Printf("\nCreated %s\n", outputFile)
	fmt.Println("Run: cvwonder generate")
	return nil
}
