import React from 'react';
import { createRoot } from 'react-dom/client';
import { Button } from './components/Button';

const App = () => (
  <div style={{ padding: 40 }}>
    <h2>Vuno UI Antd Button Demo</h2>
    <Button type="primary">Primary Button</Button>
    <Button type="default" style={{ marginLeft: 8 }}>
      Secondary Button
    </Button>
  </div>
);

const root = createRoot(document.getElementById('root')!);
root.render(<App />);
