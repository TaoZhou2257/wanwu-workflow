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

import { inject, injectable } from 'inversify';
import { EntityManager } from '@flowgram-adapter/free-layout-editor';
import {
  WorkflowDocumentProvider,
  type WorkflowDocument,
  type WorkflowNodeRegistry,
} from '@flowgram-adapter/free-layout-editor';
import {
  WorkflowBatchService,
  WorkflowVariableService,
  WorkflowVariableValidationService,
} from '@coze-workflow/variable';
import {
  type NodeTemplateInfo,
  WorkflowNodesService,
  type PlaygroundContext,
} from '@coze-workflow/nodes';
import { NODE_ORDER, StandardNodeType } from '@coze-workflow/base/types';
import {
  workflowApi,
  WorkflowMode,
  type NodeTemplateListResponse,
} from '@coze-workflow/base/api';
import { REPORT_EVENTS } from '@coze-arch/report-events';
import { CustomError } from '@coze-arch/bot-error';
import {
  type NodeCategory as ServerNodeCategory,
  type PluginCategory,
  type PluginAPINode,
} from '@coze-arch/bot-api/workflow_api';
import {
  ProductEntityType,
  SortType,
  type GetUserFavoriteListData,
} from '@coze-arch/bot-api/product_api';
import { type Model } from '@coze-arch/bot-api/developer_api';
import { ProductApi } from '@coze-arch/bot-api';

import {
  type NodeTemplate,
  type PluginApiNodeTemplate,
  type PluginCategoryNodeTemplate,
  type NodeCategory,
} from './typing';
import { createApiNodeInfo } from './hooks/use-add-node-modal/helper';
import { WorkflowGlobalStateEntity } from './entities';
import { PAGE_SIZE } from './components/node-panel/constant';

export interface ImageflowNode {
  title: string;
  desc: string;
  icon: string;
  apiName: string;
  apiID: string;
  pluginID: string;
  pluginName: string;
}

export interface ImageflowNodesGroup {
  category: string;
  tools: ImageflowNode[];
}
@injectable()
export class WorkflowPlaygroundContext implements PlaygroundContext {
  protected nodeTemplateMap = new Map<StandardNodeType, NodeTemplate>();
  protected pluginApiMap: Record<string, PluginAPINode> = {};
  protected pluginCategoryMap: Record<string, PluginCategory> = {};
  public favoritePlugins: GetUserFavoriteListData | undefined;

  protected nodeCategoryList: ServerNodeCategory[] = [];
  public imageflowNodes: ImageflowNode[];

  @inject(WorkflowDocumentProvider)
  protected documentProvider: WorkflowDocumentProvider;

  @inject(WorkflowVariableService)
  readonly variableService: WorkflowVariableService;
  @inject(WorkflowBatchService) readonly batchService: WorkflowBatchService;
  @inject(WorkflowVariableValidationService)
  readonly variableValidationService: WorkflowVariableValidationService;
  @inject(WorkflowNodesService) readonly nodesService: WorkflowNodesService;
  @inject(EntityManager) public entityManager: EntityManager;

  protected modelList: Model[];

  get models() {
    return this.modelList;
  }

  set models(models: Model[]) {
    this.modelList = models;
  }

  /**
   * Acquire documents
   */
  get document(): WorkflowDocument {
    return this.documentProvider();
  }
  /**
   * Get, workflow node template
   */
  async loadNodeInfos(locale: string): Promise<void> {
    const nodeIds: StandardNodeType[] = Object.values(StandardNodeType);
    let resp: NodeTemplateListResponse | undefined;
    let favoritePlugins: GetUserFavoriteListData | undefined;
    const response = await Promise.allSettled([
      workflowApi.NodeTemplateList(
        {
          node_types: nodeIds,
        },
        {
          headers: {
            'x-locale': locale, // zh-CN, en-US
          },
        },
      ),
      this.fetchFavoritePlugins({ pageNum: 1 }),
    ]);
    response[0].status === 'fulfilled' && (resp = {
      data: {
        "template_list": [
          {
            "id": "1",
            "type": 1,
            "name": "开始",
            "desc": "工作流的起始节点，用于设定启动工作流需要的信息",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start-v2.jpg",
            "support_batch": 1,
            "node_type": "1",
            "color": "#5C62FF"
          },
          {
            "id": "2",
            "type": 2,
            "name": "结束",
            "desc": "工作流的最终节点，用于返回工作流运行后的结果信息",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End-v2.jpg",
            "support_batch": 1,
            "node_type": "2",
            "color": "#5C62FF"
          },
          {
            "id": "13",
            "type": 13,
            "name": "输出",
            "desc": "节点从“消息”更名为“输出”，支持中间过程的消息输出，支持流式和非流式两种方式",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Output-v2.jpg",
            "support_batch": 1,
            "node_type": "13",
            "color": "#5C62FF"
          },
          {
            "id": "30",
            "type": 30,
            "name": "输入",
            "desc": "支持中间过程的信息输入",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Input-v2.jpg",
            "support_batch": 1,
            "node_type": "30",
            "color": "#5C62FF"
          },
          {
            "id": "3",
            "type": 3,
            "name": "大模型",
            "desc": "调用大语言模型,使用变量和提示词生成回复",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-LLM-v2.jpg",
            "support_batch": 2,
            "node_type": "3",
            "color": "#5C62FF"
          },
          {
            "id": "4",
            "type": 4,
            "name": "插件",
            "desc": "通过添加工具访问实时数据和执行外部操作",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Plugin-v2.jpg",
            "support_batch": 2,
            "node_type": "4",
            "color": "#CA61FF"
          },
          {
            "id": "9",
            "type": 9,
            "name": "工作流",
            "desc": "集成已发布工作流，可以执行嵌套子任务",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Workflow-v2.jpg",
            "support_batch": 2,
            "node_type": "9",
            "color": "#00B83E"
          },
          {
            "id": "5",
            "type": 5,
            "name": "代码",
            "desc": "编写代码，处理输入变量来生成返回值",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Code-v2.jpg",
            "support_batch": 1,
            "node_type": "5",
            "color": "#00B2B2"
          },
          {
            "id": "8",
            "type": 8,
            "name": "选择器",
            "desc": "连接多个下游分支，若设定的条件成立则仅运行对应的分支，若均不成立则只运行“否则”分支",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Condition-v2.jpg",
            "support_batch": 1,
            "node_type": "8",
            "color": "#00B2B2"
          },
          {
            "id": "19",
            "type": 19,
            "name": "终止循环",
            "desc": "用于立即终止当前所在的循环，跳出循环体",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Break-v2.jpg",
            "support_batch": 1,
            "node_type": "19",
            "color": "#00B2B2"
          },
          {
            "id": "20",
            "type": 20,
            "name": "设置变量",
            "desc": "用于重置循环变量的值，使其下次循环使用重置后的值",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-LoopSetVariable-v2.jpg",
            "support_batch": 1,
            "node_type": "20",
            "color": "#00B2B2"
          },
          {
            "id": "21",
            "type": 21,
            "name": "循环",
            "desc": "用于通过设定循环次数和逻辑，重复执行一系列任务",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Loop-v2.jpg",
            "support_batch": 1,
            "node_type": "21",
            "color": "#00B2B2"
          },
          {
            "id": "22",
            "type": 22,
            "name": "意图识别",
            "desc": "用于用户输入的意图识别，并将其与预设意图选项进行匹配。",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Intent-v2.jpg",
            "support_batch": 1,
            "node_type": "22",
            "color": "#00B2B2"
          },
          {
            "id": "28",
            "type": 28,
            "name": "批处理",
            "desc": "通过设定批量运行次数和逻辑，运行批处理体内的任务",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Batch-v2.jpg",
            "support_batch": 1,
            "node_type": "28",
            "color": "#00B2B2"
          },
          {
            "id": "29",
            "type": 29,
            "name": "继续循环",
            "desc": "用于终止当前循环，执行下次循环",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Continue-v2.jpg",
            "support_batch": 1,
            "node_type": "29",
            "color": "#00B2B2"
          },
          {
            "id": "32",
            "type": 32,
            "name": "变量聚合",
            "desc": "对多个分支的输出进行聚合处理",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/VariableMerge-icon.jpg",
            "support_batch": 1,
            "node_type": "32",
            "color": "#00B2B2"
          },
          {
            "id": "61",
            "type": 61,
            "name": "知识库",
            "desc": "在选定的知识中,根据输入变量召回最匹配的信息,并以列表形式返回",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeQuery-v2.jpg",
            "support_batch": 1,
            "node_type": "61",
            "color": "#FF811A"
          },
          {
            "id": "6",
            "type": 6,
            "name": "知识库检索",
            "desc": "在选定的知识中,根据输入变量召回最匹配的信息,并以列表形式返回",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeQuery-v2.jpg",
            "support_batch": 1,
            "node_type": "6",
            "color": "#FF811A"
          },
          {
            "id": "27",
            "type": 27,
            "name": "知识库写入",
            "desc": "写入节点可以添加 文本类型 的知识库，仅可以添加一个知识库",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeWriting-v2.jpg",
            "support_batch": 1,
            "node_type": "27",
            "color": "#FF811A"
          },
          {
            "id": "40",
            "type": 40,
            "name": "变量赋值",
            "desc": "用于给支持写入的变量赋值，包括应用变量、用户变量",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/Variable.jpg",
            "support_batch": 1,
            "node_type": "40",
            "color": "#FF811A"
          },
          {
            "id": "12",
            "type": 12,
            "name": "SQL自定义",
            "desc": "基于用户自定义的 SQL 完成对数据库的增删改查操作",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Database-v2.jpg",
            "support_batch": 1,
            "node_type": "12",
            "color": "#FF811A"
          },
          {
            "id": "42",
            "type": 42,
            "name": "更新数据",
            "desc": "修改表中已存在的数据记录，用户指定更新条件和内容来更新数据",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-database-update.jpg",
            "support_batch": 1,
            "node_type": "42",
            "color": "#F2B600"
          },
          {
            "id": "43",
            "type": 43,
            "name": "查询数据",
            "desc": "从表获取数据，用户可定义查询条件、选择列等，输出符合条件的数据",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icaon-database-select.jpg",
            "support_batch": 1,
            "node_type": "43",
            "color": "#F2B600"
          },
          {
            "id": "44",
            "type": 44,
            "name": "删除数据",
            "desc": "从表中删除数据记录，用户指定删除条件来删除符合条件的记录",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-database-delete.jpg",
            "support_batch": 1,
            "node_type": "44",
            "color": "#F2B600"
          },
          {
            "id": "46",
            "type": 46,
            "name": "新增数据",
            "desc": "向表添加新数据记录，用户输入数据内容后插入数据库",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-database-insert.jpg",
            "support_batch": 1,
            "node_type": "46",
            "color": "#F2B600"
          },
          {
            "id": "15",
            "type": 15,
            "name": "文本处理",
            "desc": "用于处理多个字符串类型变量的格式",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-StrConcat-v2.jpg",
            "support_batch": 1,
            "node_type": "15",
            "color": "#3071F2"
          },
          {
            "id": "18",
            "type": 18,
            "name": "问答",
            "desc": "支持中间向用户提问问题,支持预置选项提问和开放式问题提问两种方式",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Direct-Question-v2.jpg",
            "support_batch": 1,
            "node_type": "18",
            "color": "#3071F2"
          },
          {
            "id": "45",
            "type": 45,
            "name": "HTTP 请求",
            "desc": "用于发送API请求，从接口返回数据",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-HTTP.png",
            "support_batch": 1,
            "node_type": "45",
            "color": "#3071F2"
          },
          {
            "id": "62",
            "type": 62,
            "name": "MCP",
            "desc": "MCP 描述",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-to_json.png",
            "support_batch": 1,
            "node_type": "62",
            "color": "F2B600"
          },
          {
            "id": "63",
            "type": 63,
            "name": "文档解析",
            "desc": "文档解析描述",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-to_json.png",
            "support_batch": 1,
            "node_type": "63",
            "color": "F2B600"
          },
          {
            "id": "58",
            "type": 58,
            "name": "JSON 序列化",
            "desc": "用于把变量转化为JSON字符串",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-to_json.png",
            "support_batch": 1,
            "node_type": "58",
            "color": "F2B600"
          },
          {
            "id": "59",
            "type": 59,
            "name": "JSON 反序列化",
            "desc": "用于将JSON字符串解析为变量",
            "icon_url": "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-from_json.png",
            "support_batch": 1,
            "node_type": "59",
            "color": "F2B600"
          }
        ],
        "cate_list": [
          {
            "name": "",
            "node_type_list": [
              "3",
              "4",
              "9",
              "61",
              "62",
            ]
          },
          {
            "name": "业务逻辑",
            "node_type_list": [
              "5",
              "8",
              "19",
              "20",
              "21",
              "22",
              "28",
              "29",
              "32"
            ]
          },
          {
            "name": "输入\u0026输出",
            "node_type_list": [
              "1",
              "2",
              "13",
              "30"
            ]
          },
          {
            "name": "数据库",
            "node_type_list": [
              "12",
              "42",
              "43",
              "44",
              "46"
            ]
          },
          {
            "name": "知识库\u0026数据",
            "node_type_list": [
              "6",
              "27",
              "40"
            ]
          },
          {
            "name": "组件",
            "node_type_list": [
              "15",
              "18",
              "45",
              "58",
              "59",
              "63"
            ]
          }
        ],
        "plugin_api_list": null,
        "plugin_category_list": null
      }
    }) //response[0].value);
    response[1].status === 'fulfilled' && (favoritePlugins = response[1].value);

    // Convert the template data returned by the server level to type: '1', which is the same as the standard StandardNodeType.
    const typeKey = 'node_type';

    this.favoritePlugins = favoritePlugins;
    this.nodeCategoryList = resp?.data?.cate_list ?? [];
    this.pluginApiMap = (resp?.data?.plugin_api_list ?? []).reduce<
      Record<string, PluginAPINode>
    >((acc, curr) => {
      curr.api_id && (acc[curr.api_id] = curr);
      return acc;
    }, {});
    this.pluginCategoryMap = (resp?.data?.plugin_category_list ?? []).reduce<
      Record<string, PluginCategory>
    >((acc, curr) => {
      curr.plugin_category_id && (acc[curr.plugin_category_id] = curr);
      return acc;
    }, {});
    resp?.data?.template_list?.forEach(temp => {
      if (temp[typeKey]) {
        this.nodeTemplateMap.set(`${temp[typeKey]}` as StandardNodeType, {
          ...temp,
          type: `${temp[typeKey]}` as StandardNodeType,
        });
      }
    });
    console.log(this.pluginCategoryMap, this.nodeTemplateMap, this.nodeCategoryList, '-----------------------------loadNodeInfos')
  }

  getImageFlowNode(pluginId: string, apiName: string) {
    return this.imageflowNodes.find(
      tool => tool.pluginID === pluginId && tool.apiName === apiName,
    );
  }

  getNodeTemplateInfoByType = (
    type: StandardNodeType,
  ): NodeTemplateInfo | undefined => {
    const registry = this.document.getNodeRegister<WorkflowNodeRegistry>(type);
    if (
      !registry ||
      !registry.meta ||
      registry.meta.nodeDTOType === undefined
    ) {
      throw new CustomError(
        REPORT_EVENTS.parmasValidation,
        `Unknown NodeMeta by type ${type}`,
      );
    }
    const info = this.nodeTemplateMap.get(
      registry.meta.nodeDTOType as StandardNodeType,
    );

    if (!info) {
      return;
    }

    return {
      title: info.name as string,
      icon: info.icon_url as string,
      description: info.desc as string,
      mainColor: info.color || '',
      subTitle: (type !== StandardNodeType.Start &&
      type !== StandardNodeType.End
        ? info.name
        : '') as string,
    };
  };

  /**
   * This will prohibit operations such as circling and deleting
   */
  get disabled(): boolean {
    return !!this.globalState?.readonly;
  }

  /**
   * This way of passing through the context is not very good, and it needs to be rebuilt later.
   */
  get spaceId(): string | undefined {
    return this.globalState?.spaceId;
  }

  get flowMode(): WorkflowMode {
    return this.globalState?.flowMode ?? WorkflowMode.Workflow;
  }

  getTemplateList(types: StandardNodeType[] = []): NodeTemplate[] {
    // HACK: The type of the passed type is string, and then the template returned by the backend is actually number.
    return (
      types
        .sort(
          (prev, next) => Number(NODE_ORDER[prev]) - Number(NODE_ORDER[next]),
        )
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        .map(type => this.nodeTemplateMap.get(type)!)
        .filter(Boolean)
    );
  }

  async fetchFavoritePlugins({
    pageNum,
    pageSize = PAGE_SIZE,
  }: {
    pageNum: number;
    pageSize?: number;
  }): Promise<GetUserFavoriteListData | undefined> {
    // will support soon
    if (IS_OPEN_SOURCE) {
      return {
        favorite_products: [],
        has_more: false,
      };
    }

    const resp = await ProductApi.PublicGetUserFavoriteList({
      entity_type: ProductEntityType.Plugin,
      sort_type: SortType.Newest,
      page_num: pageNum,
      page_size: pageSize,
    });
    return resp.data;
  }

  getTemplateCategoryList(
    enabledNodeTypes: StandardNodeType[] = [],
    isSupportImageflowNodes = false,
  ): NodeCategory[] {
    const isBindDouyin = Boolean(this.globalState?.isBindDouyin);
    const nodeCategoryList =
      this.nodeCategoryList.length !== 0
        ? this.nodeCategoryList
        : // fallback when server data load failed
          [
            {
              name: '',
              node_type_list: enabledNodeTypes,
            },
          ];
    return nodeCategoryList
      .map(category => {
        const nodeList =
          category.node_type_list
            ?.filter(item =>
              enabledNodeTypes.includes(item as StandardNodeType),
            )
            ?.map(item => this.nodeTemplateMap.get(item as StandardNodeType))
            ?.filter<NodeTemplate>((item): item is NodeTemplate =>
              Boolean(item),
            ) ?? [];
        const pluginApiList: PluginApiNodeTemplate[] = (
          category.plugin_api_id_list ?? []
        )?.map(apiId => {
          const pluginInfo = this.pluginApiMap[apiId];
          const nodeJSON = createApiNodeInfo(
            {
              name: pluginInfo.api_name,
              plugin_name: pluginInfo.name,
              api_id: pluginInfo.api_id,
              plugin_id: pluginInfo.plugin_id,
              desc: pluginInfo.desc,
            },
            pluginInfo.icon_url,
          );
          return {
            type: StandardNodeType.Api,
            ...pluginInfo,
            nodeJSON,
          };
        });
        const pluginCategoryList: PluginCategoryNodeTemplate[] = (
          category.plugin_category_id_list ?? []
        ).map(categoryId => {
          const pluginCategory = this.pluginCategoryMap[categoryId];
          return {
            type: StandardNodeType.Api,
            ...pluginCategory,
            categoryInfo: {
              categoryId,
              onlyOfficial: pluginCategory.only_official,
            },
          };
        });
        return {
          categoryName: category.name,
          nodeList: isSupportImageflowNodes
            ? [
                ...nodeList,
                ...(isBindDouyin ? [] : pluginApiList),
                ...(isBindDouyin ? [] : pluginCategoryList),
              ]
            : nodeList,
        };
      })
      .filter(category => category.nodeList.length > 0);
  }
  get globalState() {
    return this.entityManager.getEntity<WorkflowGlobalStateEntity>(
      WorkflowGlobalStateEntity,
    );
  }

  get schemaGray(): {
    isBatchV2: boolean;
  } {
    // const { schemaGray } = this.globalState?.config ?? {};
    return {
      isBatchV2: false,
    };
  }
}
