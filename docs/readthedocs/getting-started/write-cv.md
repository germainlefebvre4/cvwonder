# Write your CV

Write your CV in YAML format following the schema designed in the json schema file.

!!! note

    **None of the fields are required.**
    You can fill in only the fields you want. It will all depend on the information printed by the theme.

## Sections

The CV is divided into sections. Each section has a title and a list of items. Each item has a title and a description.

Here is the list of sections.

### Company

The company section is used to describe the company you are working for.

| Field | Description |
| --- | --- |
| Name | The name of the company. |
| Logo | The logo of the company. |

??? example "Example"

    Add your company information in the company section:

    ```yaml
    company:
      name: zatsit
      logo: images/zatsit-logo.webp
    ```

### Person

The person section is used to describe "you".

| Field | Description |
| --- | --- |
| Name | Your name. |
| Depiction | A short description of you. |
| Profession | Your profession. |
| Location | Your location. |
| Citizenship | Your citizenship. |
| Email | Your email. |
| Site | Your site. |

??? example "Example"

    Add your person information in the person section:

    ```yaml
    person:
      name: Germain Lefebvre
      depiction: Software Engineer
      profession: Software Engineer
      location: Paris, France
      citizenship: French
      email: germain@lefebvre.fr
    ```

### Social Networks

The social networks section is used to describe your social networks.

| Field | Description |
| --- | --- |
| Github | Your Github profile. |
| Stackoverflow | Your Stackoverflow profile. |
| Linkedin | Your Linkedin profile. |
| Twitter | Your Twitter profile. |

??? example "Example"

    Add your social networks in the social networks section:

    ```yaml
    socialNetworks:
      github: germainlefebvre4
      stackoverflow: germainlefebvre4
      linkedin: germainlefebvre4
      twitter: germainlefebvr4
    ```

### Abstract

The abstract section is used to describe a short abstract about you.

| Field | Description |
| --- | --- |
| []Tr | The abstract text. |

You can add multiple abstracts.

??? example "Example"

    Add 2 lines in the abstract section:

    ```yaml
    abstract:
      - tr: "I am a software engineer with 5 years of experience in the field."
      - tr: "I am passionate about technology and I love to learn new things."
    ```

### Career

The career section is used to describe your multiple career paths and experiences.

| Field | Description |
| --- | --- |
| CompanyName | The name of the company. |
| CompanyLogo | The logo of the company. |
| Duration | The duration of your experience. |
| Missions | The missions you have accomplished. |

??? example "Example"

    Add your career information in the career section:

    ```yaml
    career:
      - companyName: zatsit
        companyLogo: images/zatsit-logo.webp
        duration: 2 years
        missions:
          - position: Software Engineer
            company: zatsit
            location: Paris, France
            dates: 2020-2022
            summary: Developed a web application using React and Node.js.
            description:
              - Developed the front-end using React.
              - Developed the back-end using Node.js.
            technologies:
              - React
              - Node.js
            project: MyProject
    ```

#### Missions

The missions sub-section is used to describe your missions in detail.

| Field | Description |
| --- | --- |
| Position | Your position/title in this mission. |
| Company | The company you worked for. |
| Location | The location of the mission/company. |
| Dates | The dates of your experience. |
| Summary | A summary, as short description, of your experience. |
| Description | A full description of your experience. |
| Technologies | The technologies you used in this mission. |
| Project | The name of the project you worked on. |

!!! abstract "Date format"
    The `Dates` fields is a text field. You can use any format you want. The theme will not parse it.

??? example "Example"

    Add your career information in the career section:

    ```yaml
    career:
      - companyName: ...
        missions:
          - position: Software Engineer
            company: zatsit
            location: Paris, France
            dates: 2020-2022
            summary: Developed a web application using React and Node.js.
            description:
              - Developed the front-end using React.
              - Developed the back-end using Node.js.
            technologies:
              - React
              - Node.js
            project: MyProject
    ```

### Technical Skills

The technical skills section is used to describe your technical skills.

| Field | Description |
| --- | --- |
| Domains | The domains you are skilled in. |

#### Domains

The domains sub-section is used to describe your domains and competencies.

| Field | Description |
| --- | --- |
| Name | The name of the domain. |
| Competencies | The competencies you have in this domain. |

??? example "Example"

    Add your technical skills in the technical skills section:

    ```yaml
    technicalSkills:
      domains:
        - name: Front-end
          competencies:
            - name: React
              level: 5
            - name: Angular
              level: 3
        - name: Back-end
          competencies:
            - name: Node.js
              level: 4
    ```

!!! abstract "Level"
    The `Level` field is a number from 0 to 100. It represents your level of expertise in this competency.

### Side Projects

The side projects section is used to describe the personal projects or open-source projects you have worked on.

| Field | Description |
| --- | --- |
| Name | The name of the side project. |
| Position | Your position/title in this side project. |
| Description | A short description of the side project. |
| Link | The link to the side project. |
| Type | The type of the side project. |
| Langs | The languages used in the side project. |
| Color | The color of the side project. |

??? example "Example"

    Add your side projects in the side projects section:

    ```yaml
    sideProjects:
      - name: MyProject
        position: Software Engineer
        description: A web application using React and Node.js.
        link: https://germainlefebvre.fr/myproject
        type: github
        langs: React, Node.js
        color: 61dafb
    ```

The `Type` field can be used to identify the hosting platform of the side project. The theme can use this information to display the right icon.

### Certifications

The certifications section is used to describe the certifications you have.

| Field | Description |
| --- | --- |
| CompanyName | The name of the company that issued the certification. |
| CertificationName | The name of the certification. |
| Issuer | The issuer of the certification. |
| Date | The date of the certification. |
| Link | The link to the certification. |
| Badge | The badge of the certification. |

??? example "Example"

    Add your certifications in the certifications section:

    ```yaml
    certifications:
      - companyName: AWS
          certificationName: Solutions Architect Associate
          issuer: Coursera
          date: Mars 2018
          link: https://www.credly.com/badges/dd09dc40-9ef8-43a4-addb-d861d4dadf26/public_url
          badge: images/aws-certified-solutions-architect-associate.png
    ```

### Education

The education section is used to describe your school and education.

| Field | Description |
| --- | --- |
| SchoolName | The name of the school. |
| SchoolLogo | The logo of the school. |
| Degree | The degree you obtained. |
| Location | The location of the school. |
| Dates | The dates of your education. |
| Link | The link to the school. |

??? example "Example"

    Add your education in the education section:

    ```yaml
    education:
      - schoolName: University of Paris
        schoolLogo: images/university-of-paris-logo.webp
        degree: Master's degree in Computer Science
        location: Paris, France
        dates: 2015-2018
        link: https://www.univ-paris.fr/
    ```

!!! abstract "Date format"

    The `Dates` field is a text field. You can use any format you want. The theme will not parse it.

## Full Example

Here is a full example in a `cv.yml` file:

```yaml
{!getting-started/cv.yml!}
```
