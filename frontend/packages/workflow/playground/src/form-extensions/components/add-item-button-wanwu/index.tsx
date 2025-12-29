/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { type FC } from 'react';

import classnames from 'classnames';
import { IconCozPlus } from '@coze-arch/coze-design/icons';
import { Button, Tooltip } from '@coze-arch/coze-design';
import { I18n } from '@coze-arch/i18n';

interface Props {
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
  onMouseDown?: React.MouseEventHandler<HTMLButtonElement>;
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  disabledTooltip?: string;
  testId?: string;
  isAdded?: boolean;
  visible?: boolean;
}

export const AddItemWanwuButton: FC<Props> = ({
  onClick,
  onMouseDown,
  disabled,
  disabledTooltip,
  className,
  style,
  testId,
  isAdded = false,
  visible = true,
}) => {
  if (!visible) {
    return null;
  }

  const buttonContent = isAdded ? (
    <Button
      color="secondary"
      size="small"
      disabled={true}
      className={classnames('!block', className)}
      style={style}
    >
      {I18n.t('Added')}
    </Button>
  ) : (
    <Button
      color="highlight"
      size="small"
      icon={<IconCozPlus className="text-sm" />}
      onMouseDown={e => onMouseDown?.(e)}
      onClick={(e: React.MouseEvent<HTMLButtonElement>) => {
        if (disabled) {
          return;
        }
        onClick?.(e);
      }}
      className={classnames('!block', className)}
      style={style}
      disabled={disabled}
      data-testid={testId}
    >
      {I18n.t('Add_1')}
    </Button>
  );

  if (isAdded) {
    return buttonContent;
  }

  if (disabled) {
    if (disabledTooltip) {
      return <Tooltip content={disabledTooltip}>{buttonContent}</Tooltip>;
    }
    return null;
  }

  return buttonContent;
};
