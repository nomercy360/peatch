import type { UserPreview } from './UserPreview';

/**
 * @description OK
 */
export type GetApiUsersPreview200 = UserPreview[];

/**
 * @description OK
 */
export type GetApiUsersPreviewQueryResponse = UserPreview[];

export type GetApiUsersPreviewQuery = {
  Response: GetApiUsersPreviewQueryResponse;
};