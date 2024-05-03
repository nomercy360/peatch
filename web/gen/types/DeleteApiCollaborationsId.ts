import { CreateCollaborationRequest } from './CreateCollaborationRequest';

export type DeleteApiCollaborationsIdPathParams = {
    /**
     * @description Collaboration ID
     * @type integer
    */
    id: number;
};

 /**
 * @description No Content
*/
export type DeleteApiCollaborationsId204 = any;

/**
 * @description Collaboration data
 */
export type DeleteApiCollaborationsIdMutationRequest = CreateCollaborationRequest;

 export type DeleteApiCollaborationsIdMutationResponse = any;

 export type DeleteApiCollaborationsIdMutation = {
    Response: DeleteApiCollaborationsIdMutationResponse;
   Request: DeleteApiCollaborationsIdMutationRequest;
    PathParams: DeleteApiCollaborationsIdPathParams;
};