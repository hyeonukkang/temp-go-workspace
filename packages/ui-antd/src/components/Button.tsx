import React from 'react';
import { Button as AntdButton, ButtonProps as AntdButtonProps } from 'antd';

export type ButtonProps = AntdButtonProps;

export const Button: React.FC<ButtonProps> = (props) => {
  return <AntdButton {...props} />;
};
