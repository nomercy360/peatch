import type { Collaboration } from "./Collaboration";

 export type GetApiCollaborationsIdPathParams = {
    /**
     * @description Collaboration ID
     * @type integer
    */
    id: number;
};

 /**
 * @description OK
*/
export type GetApiCollaborationsId200 = Collaboration;

 /**
 * @description OK
*/
export type GetApiCollaborationsIdQueryResponse = Collaboration;

 export type GetApiCollaborationsIdQuery = {
    Response: GetApiCollaborationsIdQueryResponse;
    PathParams: GetApiCollaborationsIdPathParams;
};