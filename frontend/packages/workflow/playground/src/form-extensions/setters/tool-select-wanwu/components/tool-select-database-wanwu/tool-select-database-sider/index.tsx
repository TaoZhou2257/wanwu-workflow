import { useState, type FC } from 'react';

import { I18n } from '@coze-arch/i18n';
import { UICompositionModalSider } from '@coze-arch/bot-semi';

export interface DatabaseModalSiderProps {
  activeTab: 'database' | 'qaBase';
  renderSearch?: () => React.ReactNode;
  renderQaSearch?: () => React.ReactNode;
  onTabChange: (key: string) => void;
}

const tabList = [
  {
    key: 'database',
    label: I18n.t('Datasets' as any, {}, '知识库'),
  },
  /*{
    key: 'qaBase',
    label: I18n.t('workflow_detail_qa_knowledge_knowledge' as any, {}, '问答库'),
  },*/
];

export const DatabaseModalSider: FC<DatabaseModalSiderProps> = ({
  renderSearch,
  renderQaSearch,
  activeTab,
  onTabChange,
}) => {
  const [activeKey, setActiveKey] = useState<string>('database');
  const searchRender = activeTab === 'qaBase' ? renderQaSearch : renderSearch;

  return (
    <UICompositionModalSider style={{ paddingTop: 16 }}>
      <UICompositionModalSider.Header>
        {searchRender ? searchRender() : null}
      </UICompositionModalSider.Header>

      <UICompositionModalSider.Content style={{ paddingTop: 16 }}>
        {tabList.map(item => {
          const isActive = item.key === activeKey;
          return (
            <div
              key={item.key}
              className={
                isActive
                  ? 'px-[12px] py-[10px] rounded-[8px] text-[14px] font-semibold mb-[8px] bg-[rgba(46,47,56,0.05)] text-[var(--Text-coz-text-primary,#1d2129)] cursor-pointer'
                  : 'px-[12px] py-[10px] rounded-[8px] text-[14px] mb-[8px] text-[var(--Text-coz-text-primary,#1d2129)] cursor-pointer'
              }
              onClick={() => {
                setActiveKey(item.key);
                onTabChange(item.key);
              }}
            >
              {item.label}
            </div>
          );
        })}
      </UICompositionModalSider.Content>
    </UICompositionModalSider>
  );
};

