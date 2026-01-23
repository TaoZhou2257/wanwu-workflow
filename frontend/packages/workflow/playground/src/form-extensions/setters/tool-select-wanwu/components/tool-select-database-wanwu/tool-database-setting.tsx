import { useEffect, useState } from 'react';
import { Modal, Button, Toast } from '@coze-arch/coze-design';
import { I18n } from '@coze-arch/i18n';
import { isEmpty } from 'lodash-es';

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
    matchType,
    rerankModelId = '',
  } = value || {};

  const [isInit, setIsInit] = useState<boolean>(true);
  const { dataSets, isReady } = useDataSetInfos({ 
    ids: selectDataSet.map(item => item?.dataset_id || item?.id).filter(Boolean) 
  });

  const needRerankIdList:any = [MatchType.Semantic, MatchType.FullText, MatchType.Hybird]
  const isAllExternalKnowledge = selectDataSet.every(item => item.external)

  useEffect(() => {
    if (!isInit && !isEmpty(value)) {
      if (isAllExternalKnowledge) {
        onChange?.({
          ...value,
          useGraph: false,
          rewrite: false,
        });
      }
    }
    setIsInit(false);
  }, [isAllExternalKnowledge]);

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
            if (!isAllExternalKnowledge && needRerankIdList.includes(matchType) && !rerankModelId) {
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

