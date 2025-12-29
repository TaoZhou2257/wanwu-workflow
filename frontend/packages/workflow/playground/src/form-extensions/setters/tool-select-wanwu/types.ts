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
export interface Tool {
  actionID: string;
  actionName: string;
  toolType: string;
  apiKey: string;
  toolId: string;
  toolName: string;
}

export type ToolSelectValue = Tool[];

export interface DatabaseSelectContextProps {
  changeTool: (id: string) => void;
  clearTool: () => void;
  readonly?: boolean;
}

export const TOOL_TAB = {
  TOOL: 'tool',
  WORKFLOW: 'workflow',
  DATABASE: 'database',
  MCP: 'mcp',
} as const;
