import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import HomepageDemo from '@site/src/components/HomepageDemo';
import HomepageThemes from '@site/src/components/HomepageThemes';
import Heading from '@theme/Heading';
import IconCVWonderLogo from '@site/static/img/logo.svg';

import styles from './index.module.css';

function HomepageLogo() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.logoBanner)}>
      <div className="container" style={{display: 'grid', gridTemplateColumns: '', justifyContent: '',}}>
        <div style={{minWidth: '400px',}}>
          <IconCVWonderLogo style={{maxWidth: '400px',maxHeight: '400px',}} />
        </div>
      </div>
    </header>
  );
}

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container" style={{display: 'grid', gridTemplateColumns: '', justifyContent: '',}}>
        <div style={{minWidth: '400px',}}>
          {/* <IconCVWonderLogo style={{maxWidth: '400px',maxHeight: '400px',}} /> */}
          <Heading as="h1" className="hero__title">
            {siteConfig.title}
          </Heading>
          <p className="hero__subtitle">{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className="button button--secondary button--lg margin-right--md"
              to="/docs/getting-started">
              Getting Started
            </Link>
            <Link
              className="button button--outline button--lg button--secondary"
              href="https://github.com/germainlefebvre4/cvwonder">
              View on GitHub
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
}

function HomepageHeaderBackup() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg margin-right--md"
            to="/docs/getting-started">
            Getting Started
          </Link>
          <Link
            className="button button--outline button--lg button--secondary"
            href="https://github.com/germainlefebvre4/cvwonder">
            View on GitHub
          </Link>
        </div>
      </div>
    </header>
  );
}

function HomepageCTA() {
  return (
    <section className={styles.ctaSection}>
      <div className="container text--center">
        <Heading as="h2">Ready to Create Your CV?</Heading>
        <p className={styles.ctaText}>
          Get started with CV Wonder today and create beautiful, professional CVs in minutes
        </p>
        <div className={styles.ctaButtons}>
          <Link
            className="button button--primary button--lg margin-right--md"
            to="/docs/getting-started">
            Installation Guide
          </Link>
          <Link
            className="button button--secondary button--lg"
            to="/docs/cli">
            CLI Documentation
          </Link>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - Professional CV Generator`}
      description="Generate beautiful, professional CVs from YAML in seconds with CV Wonder. Choose from multiple themes or create your own.">
      <HomepageLogo />
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <HomepageDemo />
        <HomepageThemes />
        <HomepageCTA />
      </main>
    </Layout>
  );
}
