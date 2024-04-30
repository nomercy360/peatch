import { CreateUserCollaboration } from "./CreateUserCollaboration";
import type { Collaboration } from "./Collaboration";

 export type PostApiUsersUserIdCollaborationsPathParams = {
    /**
     * @description User ID
     * @type integer
    */
    user_id: number;
};

 /**
 * @description OK
*/
export type PostApiUsersUserIdCollaborations200 = Collaboration;

 /**
 * @description Collaboration data
*/
export type PostApiUsersUserIdCollaborationsMutationRequest = CreateUserCollaboration;

 /**
 * @description OK
*/
export type PostApiUsersUserIdCollaborationsMutationResponse = Collaboration;

 export type PostApiUsersUserIdCollaborationsMutation = {
    Response: PostApiUsersUserIdCollaborationsMutationResponse;
    Request: PostApiUsersUserIdCollaborationsMutationRequest;
    PathParams: PostApiUsersUserIdCollaborationsPathParams;
};