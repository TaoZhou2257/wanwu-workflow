/*
 * 只包含 workflow 的路由
 */

import {createBrowserRouter, Navigate} from 'react-router-dom';
import { lazy } from 'react';
import { GlobalError } from '@coze-foundation/layout';

const LoginPage = lazy(() =>
    import('@coze-foundation/account-ui-adapter').then(res => ({
      default: res.LoginPage,
    })),
);
const WorkflowPage = lazy(() =>
  import('@coze-workflow/playground-adapter').then(res => ({
    default: res.WorkflowPage,
  })),
);

export const workflowRouter: ReturnType<typeof createBrowserRouter> =
  createBrowserRouter([
    {
      path: '/',
      // errorElement: <GlobalError />,
      children: [
        /*{
          path: 'sign',
          Component: LoginPage,
          // errorElement: <GlobalError />,
          loader: () => ({
            hasSider: false,
            requireAuth: false,
          }),
        },*/
        {
          index: true,
          element: <Navigate to="/workflow" replace />,
        },
        {
          path: 'workflow',
          Component: WorkflowPage,
          loader: () => ({
            hasSider: false,
            requireAuth: false,
          }),
        },
      ],
    },
  ]);
