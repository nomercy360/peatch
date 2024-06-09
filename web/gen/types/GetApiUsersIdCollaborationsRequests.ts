import type { CollaborationRequest } from "./CollaborationRequest";

 export type GetApiUsersIdCollaborationsRequestsPathParams = {
    /**
     * @description User ID
     * @type integer
    */
    id: number;
};

 /**
 * @description OK
*/
export type GetApiUsersIdCollaborationsRequests200 = CollaborationRequest;

 /**
 * @description OK
*/
export type GetApiUsersIdCollaborationsRequestsQueryResponse = CollaborationRequest;

 export type GetApiUsersIdCollaborationsRequestsQuery = {
    Response: GetApiUsersIdCollaborationsRequestsQueryResponse;
    PathParams: GetApiUsersIdCollaborationsRequestsPathParams;
};