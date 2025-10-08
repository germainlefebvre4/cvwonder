import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import Heading from '@theme/Heading';
import Link from '@docusaurus/Link';

const themes = [
  {
    title: 'Default Theme',
    image: require('@site/static/img/theme-default-w800px.png').default,
    description: 'Clean, professional design suitable for all industries'
  },
  {
    title: 'Horizon Timeline',
    image: require('@site/static/img/theme-horizon-timeline-w800px.png').default,
    description: 'Modern timeline layout with visual emphasis on career progression'
  },
  {
    title: 'Create Your Own',
    image: require('@site/static/img/theme-create-your-own-w800px.png').default,
    description: 'Design your own themes with our flexible templating system'
  }
];

export default function HomepageThemes() {
  return (
    <section className={clsx(styles.themesSection, 'bg-alt')}>
      <div className="container">
        <div className="text--center margin-bottom--xl">
          <Heading as="h2">Beautiful Themes For Every Need</Heading>
          <p className={styles.sectionSubtitle}>
            Choose from our collection of professionally designed themes or create your own
          </p>
        </div>
        <div className="row">
          {themes.map((theme, idx) => (
            <div className="col col--4" key={idx}>
              <div className={styles.themeCard}>
                <div className={styles.themeImageWrapper}>
                  <div className={styles.themeImagePlaceholder}>
                    <img
                      src={theme.image}
                      alt={theme.title}
                    />
                    <span>{theme.title}</span>
                    <small>Replace with theme screenshot</small>
                  </div>
                </div>
                <h3>{theme.title}</h3>
                <p>{theme.description}</p>
              </div>
            </div>
          ))}
        </div>
        <div className="text--center margin-top--lg">
          <Link
            className="button button--outline button--secondary button--lg"
            to="/docs/themes/library">
            Explore All Themes
          </Link>
        </div>
      </div>
    </section>
  );
}
