import { useState, useMemo } from 'react';

import { Modal, Popover, Button, Toast } from '@coze-arch/coze-design';
import { IconWarningInfo } from '@coze-arch/bot-icons';
import { I18n } from '@coze-arch/i18n';

import {
  MatchType,
  type DataSetInfo,
} from '../../../../components/dataset-setting/type';
import { DataSetSetting } from '../../../../components/dataset-setting/index-wanwu';
import { useDataSetInfos } from '@/hooks';
import s from '../../../../components/dataset-setting/index.module.less';

interface ToolDatabaseSettingProps {
  visible: boolean;
  value: DataSetInfo;
  onChange: (v: DataSetInfo) => void;
  onOk?: () => void;
  onCancel?: () => void;
  readonly?: boolean;
  disabled?: boolean;
  showUseGraph?: boolean;
  selectDataSet?: any[];
}

const SUGGEST_TOP_K = 5;
const DEFAULT_MIN_SCORE = 0.4;
const DEFAULT_TOP_K = 5;

export const ToolDatabaseSetting = ({
  visible,
  value,
  onChange,
  onOk,
  onCancel,
  readonly,
  disabled,
  showUseGraph = false,
  selectDataSet = [],
}: ToolDatabaseSettingProps) => {
  const {
    matchType = MatchType.FullText,
    rerankModelId = '',
    topK = DEFAULT_TOP_K,
    threshold = DEFAULT_MIN_SCORE,
    rewrite = true,
    useGraph = false,
  } = value || {};

  const { dataSets, isReady } = useDataSetInfos({ 
    ids: selectDataSet.map(item => item?.dataset_id || item?.id).filter(Boolean) 
  });

  return (
    <Modal
      title={I18n.t('knowledge_setting' as any, {}, '知识库设置')}
      visible={visible}
      onOk={onOk}
      onCancel={onCancel}
      width={480}
    >
      <DataSetSetting
        selectDataSet={selectDataSet}
        dataSetInfo={value}
        onDataSetInfoChange={onChange}
        readonly={readonly}
        disabled={disabled}
        isReady={isReady}
        dataSets={dataSets}
      />
      <div className={s['setting-item']} style={{ marginTop: '16px', display: 'flex', justifyContent: 'center', width: '100%' }}>
        <Button
          type="primary"
          onClick={() => {
            if (matchType !== 'mix_priority' && !rerankModelId) {
              Toast.error(I18n.t('knowledge_rerank_placeholder' as any, {}, '请选择rerank模型'));
              return;
            }
            onOk?.();
          }}
          disabled={readonly || disabled}
        >
          {I18n.t('confirm' as any, {}, '确定')}
        </Button>
      </div>
    </Modal>
  );
};

