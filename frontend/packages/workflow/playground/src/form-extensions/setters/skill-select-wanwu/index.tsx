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

import {
  type FC,
  useRef,
  useEffect,
  useState,
  type ReactNode,
} from 'react';

import classNames from 'classnames';
import { useInfiniteScroll, useBoolean } from 'ahooks';
import { I18n } from '@coze-arch/i18n';
import {
  Spin,
  Button,
  Toast,
} from '@coze-arch/coze-design';
import { Typography, UIEmpty } from '@coze-arch/bot-semi';
import { type ButtonProps } from '@coze-arch/bot-semi/Button';
import { workflowApi } from '@coze-workflow/base/api';

import { type SkillInfo, type SkillType } from './types';
import { SkillSelectSider } from './skill-select-sider';
interface SkillSelectProps {
  spaceId: string;
  projectID?: string;
  onAddSkill: (item: SkillInfo) => void;
  onRemoveSkill: (id: string) => void;
  skillList?: string[];
  tips?: ReactNode;
}

interface GetSkillListData {
  list: SkillInfo[];
  nextOffset: number;
  total: number;
  hasMore: boolean | undefined;
}

const AddedButton = (buttonProps: ButtonProps) => {
  const [isMouseIn, { setFalse, setTrue }] = useBoolean(false);

  const onMouseEnter = () => {
    setTrue();
  };
  const onMouseLeave = () => {
    setFalse();
  };

  return (
    <Button
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
      {...buttonProps}
      style={isMouseIn ? { color: 'red' } : {}}
      className="w-[75px]"
      color="primary"
    >
      {isMouseIn ? I18n.t('Remove') : I18n.t('Added')}
    </Button>
  );
};

// eslint-disable-next-line max-lines-per-function
export const SkillSelect: FC<SkillSelectProps> = ({
  spaceId,
  projectID,
  onAddSkill,
  onRemoveSkill,
  skillList = [],
  tips,
}) => {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [activeTab, setActiveTab] = useState<SkillType>('builtin');

  const getSkillType = (tab: SkillType) => {
    if (tab === 'all') return '';
    return tab;
  };

  const fetchSkillList = async (reqParams: {
    page_offset: number;
    name?: string;
    skillType?: string;
  }) => {
    const { page_offset, name = '', skillType } = reqParams;
    try {
      const response = await workflowApi.getSkillsWanwu({
        name,
        skillType,
      });
      const list = response?.data?.list || [];
      return {
        list,
        nextOffset: page_offset + 1,
        total: list?.length || 0,
        hasMore: false,
      };
    } catch (error) {
      const { statusText, data } = (error as any)?.response || {};
      Toast.error(data?.msg || statusText || 'Failed to fetch skill list');
      return {
        list: [],
        nextOffset: page_offset + 1,
        total: 0,
        hasMore: false,
      };
    }
  };

  const changeTab = (tab: SkillType) => {
    setActiveTab(tab);
  }

  const { loading, data, loadingMore, reload } = useInfiniteScroll(
    (newData?: GetSkillListData): Promise<GetSkillListData> =>
      fetchSkillList({
        page_offset: newData?.nextOffset || 0,
        name: searchKeyword,
        skillType: getSkillType(activeTab),
      }),
    {
      manual: true,
      isNoMore: newData => true,
      reloadDeps: [searchKeyword, activeTab],
      target: scrollRef,
    },
  );

  const handleAddSkill = (item: SkillInfo) => {
    if (onAddSkill && item.skillId) {
      onAddSkill(item);
    }
  };

  const handleRemoveSkill = (item: SkillInfo) => {
    if (onRemoveSkill && item.skillId) {
      onRemoveSkill(item.skillId);
    }
  };

  const isAdded = (skillId: string) => {
    return Boolean(skillList?.includes(skillId));
  };

  useEffect(() => {
    reload();
  }, []);

  // 技能列表
  const filteredSkillList = data?.list || [];

  const renderSkillCard = (item: SkillInfo) => (
    <div
      key={item.skillId}
      className="p-2.5 flex items-center gap-2.5 hover:bg-[var(--light-usage-fill-color-fill-0, rgba(46, 47, 56, .04))]"
      style={{
        borderBottom: '1px solid #dfdfdf',
      }}
    >
      <div className="flex-shrink-0">
        <img
          src={`/user/api${item.avatar?.path || ''}`}
          alt={item.skillName}
          className="w-[40px] h-[40px] rounded"
          onError={(e) => {
            e.currentTarget.src = '/v1/static/icon/skill-default-icon.png';
          }}
        />
      </div>
      <div className="flex-1 overflow-hidden">
        <div className="flex items-center mb-1">
          <Typography.Text className="!coz-fg-primary !leading-mini !font-medium !text-base truncate">
            {item.skillName}
          </Typography.Text>
        </div>
        <Typography.Text
          className="!coz-fg-secondary !leading-mini !font-normal !text-base truncate"
          ellipsis={{ showTooltip: false }}
        >
          {item.desc}
        </Typography.Text>
      </div>
      <div>
        {isAdded(item.skillId) ? (
          /*<Button
            className="w-[53px] flex justify-center items-center"
            color="primary"
            onClick={() => handleRemoveSkill(item)}
          >
            {I18n.t('Remove')}
          </Button>*/
          <AddedButton
            onClick={e => {
              e.stopPropagation();
              handleRemoveSkill(item);
            }}
          />
        ) : (
          <Button
            data-testid="bot.skill.add.modal.add.button"
            className="w-[75px] flex justify-center items-center"
            color="primary"
            onClick={() => handleAddSkill(item)}
          >
            {I18n.t('Add_2')}
          </Button>
        )}
      </div>
    </div>
  );

  const renderList = () => (
    <div className="w-full h-full flex flex-row min-h-0">
      <div className="w-[220px] mr-[12px] border-r border-[rgba(255,255,255,0.06)] h-full overflow-y-auto">
        <SkillSelectSider
          activeTab={activeTab}
          searchKeyword={searchKeyword}
          onTabChange={changeTab}
          onSearchChange={setSearchKeyword}
          searchPlaceholder={I18n.t('Search', {}, '搜索')}
        />
      </div>
      <div
        className="flex-1 flex flex-col h-full overflow-x-hidden"
        style={{ height: 'calc(100vh - 220px)' }}
      >
        <div
          className="overflow-y-auto relative flex-1"
          ref={scrollRef}
        >
          {filteredSkillList?.length === 0 ? (
            <div className="text-center py-8 coz-fg-secondary">
              {I18n.t('empty_text' as any, {}, '暂无数据')}
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-[12px]">
              {filteredSkillList?.map((item: SkillInfo) => renderSkillCard(item))}
            </div>
          )}
          {loadingMore ? (
            <div className="text-center py-4">
              {I18n.t('Loading')}
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );

  const renderEmpty = () => (
    <div className="w-full h-full flex flex-row min-h-0">
      <div className="w-[220px] mr-[12px] border-r border-[rgba(255,255,255,0.06)] h-full overflow-y-auto">
        <SkillSelectSider
          activeTab={activeTab}
          searchKeyword={searchKeyword}
          onTabChange={changeTab}
          onSearchChange={setSearchKeyword}
          searchPlaceholder={I18n.t('Search', {}, '搜索')}
        />
      </div>
      <div className="flex-1 flex flex-col h-full overflow-hidden">
        <div className="overflow-y-auto relative w-full flex-1 flex justify-center items-center">
          <UIEmpty
            isNotFound
            className="text-center py-8 coz-fg-secondary"
            notFound={{
              title: I18n.t('inifinit_search_not_found', {}, '没有找到内容'),
            }}
          ></UIEmpty>
        </div>
      </div>
    </div>
  );

  return (
    <>
      {tips}
      <div className="flex-1 min-h-0 h-full overflow-hidden">
        <Spin
          spinning={loading}
          wrapperClassName={classNames(['overflow-hidden'])}
          style={{ height: '100%' }}
        >
          {data?.list?.length !== 0 ? renderList() : renderEmpty()}
        </Spin>
      </div>
    </>
  );
};
