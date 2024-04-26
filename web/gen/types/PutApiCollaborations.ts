import type { Collaboration } from "./Collaboration";

 /**
 * @description OK
*/
export type PutApiCollaborations200 = Collaboration;

 /**
 * @description Collaboration data
*/
export type PutApiCollaborationsMutationRequest = Collaboration;

 /**
 * @description OK
*/
export type PutApiCollaborationsMutationResponse = Collaboration;

 export type PutApiCollaborationsMutation = {
    Response: PutApiCollaborationsMutationResponse;
    Request: PutApiCollaborationsMutationRequest;
};