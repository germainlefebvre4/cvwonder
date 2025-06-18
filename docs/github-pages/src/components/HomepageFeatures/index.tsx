import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Fast & Efficient',
    Svg: require('@site/static/img/2838939_417701-PDQC57-237.svg').default,
    description: (
      <>
        Generate beautiful CVs in seconds without any formatting headaches.
        Focus on your content, let CV Wonder handle the presentation.
      </>
    ),
  },
  {
    title: 'Themeable & Customizable',
    Svg: require('@site/static/img/4266119_2257832.svg').default,
    description: (
      <>
        Choose from community-built themes or create your own custom designs.
        The powerful theme system gives you complete control over the look and feel.
      </>
    ),
  },
  {
    title: 'Mass Generation Ready',
    Svg: require('@site/static/img/1311988_161361-OVJH8F-605.svg').default,
    description: (
      <>
        Need to create CVs for thousands of people? CV Wonder can handle it without breaking a sweat.
        Perfect for HR departments and recruitment agencies.
      </>
    ),
  },
];

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
