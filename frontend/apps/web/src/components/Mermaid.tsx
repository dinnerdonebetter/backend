import { FC, useEffect, useRef } from 'react';
import mermaid from 'mermaid';
import { useMantineColorScheme } from '@mantine/core';

export interface MermaidProps {
  chartDefinition: string;
}

export const Mermaid: FC<MermaidProps> = ({ chartDefinition }) => {
  const ref = useRef<HTMLDivElement>(null);

  const { colorScheme } = useMantineColorScheme();

  mermaid.mermaidAPI.initialize({
    startOnLoad: true,
    securityLevel: 'loose',
    logLevel: 5,

    theme: colorScheme === 'dark' ? 'dark' : 'base',
    themeVariables: {
      darkMode: colorScheme === 'dark',
      primaryColor: '#FFF',
      primaryTextColor: '#666',
      noteTextColor: '#666',
      primaryBorderColor: '#000',
      noteBorderColor: '#000',
      lineColor: 'rebeccapurple',
      secondaryColor: '#FFF',
    },
    darkMode: colorScheme === 'dark',
    fontFamily: 'Outfit, Helvetica, Arial, Verdana, Tahoma, Trebuchet MS, sans-serif',
    wrap: false,
  });

  useEffect(() => {
    mermaid.render('graphDiv', chartDefinition).then((result) => {
      if (ref.current) {
        ref.current.innerHTML = result.svg;
      }
    });
  }, [chartDefinition]);

  setTimeout(() => {
    mermaid.render('graphDiv', chartDefinition).then((result) => {
      if (ref.current) {
        ref.current.innerHTML = result.svg;
      }
    });
  }, 100);

  return <div key="chart" ref={ref} />;
};
