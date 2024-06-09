import type { CollaborationRequest } from "./CollaborationRequest";

 export type GetApiCollaborationsIdRequestsPathParams = {
    /**
     * @description Collaboration ID
     * @type integer
    */
    id: number;
};

 /**
 * @description OK
*/
export type GetApiCollaborationsIdRequests200 = CollaborationRequest;

 /**
 * @description OK
*/
export type GetApiCollaborationsIdRequestsQueryResponse = CollaborationRequest;

 export type GetApiCollaborationsIdRequestsQuery = {
    Response: GetApiCollaborationsIdRequestsQueryResponse;
    PathParams: GetApiCollaborationsIdRequestsPathParams;
};