import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import Heading from '@theme/Heading';
import Link from '@docusaurus/Link';

export default function HomepageDemo() {
  return (
    <section className={styles.demoSection}>
      <div className="container">
        <div className={clsx('row', styles.demoRow)}>
          <div className={clsx('col col--6', styles.demoTextCol)}>
            <Heading as="h2">Forget About Formatting</Heading>
            <p>
              Don't waste any more time formatting your CV. With CV Wonder, you can focus entirely on your content
              while our powerful engine handles all the layout and styling details. From beautiful typography to
              perfect spacing, we've got you covered.
            </p>
            <p>
              Simply provide your information in a simple YAML format, choose a theme, and let CV Wonder transform
              it into a professional CV in seconds.
            </p>
            <div className={styles.demoButtons}>
              <Link
                className="button button--primary margin-right--md"
                to="/docs/export">
                Learn About CV Format
              </Link>
            </div>
          </div>
          <div className={clsx('col col--6', styles.demoImageCol)}>
            <div className={styles.cvPreviewWrapper}>
              {/* Replace this with an actual screenshot of a CV generated with CV Wonder */}
              <div className={styles.cvPreviewPlaceholder}>
                <img
                  src={require('@site/static/img/theme-default-w800px.png').default}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
