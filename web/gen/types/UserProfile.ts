import type { Badge } from "./Badge";
import type { Opportunity } from "./Opportunity";

 export type UserProfile = {
    /**
     * @type string | undefined
    */
    avatar_url?: string;
    /**
     * @type array | undefined
    */
    badges?: Badge[];
    /**
     * @type string | undefined
    */
    city?: string;
    /**
     * @type string | undefined
    */
    country?: string;
    /**
     * @type string | undefined
    */
    country_code?: string;
    /**
     * @type string | undefined
    */
    created_at?: string;
    /**
     * @type string | undefined
    */
    description?: string;
    /**
     * @type string | undefined
    */
    first_name?: string;
    /**
     * @type integer | undefined
    */
    id?: number;
    /**
     * @type boolean | undefined
    */
    is_following?: boolean;
    /**
     * @type string | undefined
    */
    last_name?: string;
    /**
     * @type array | undefined
    */
    opportunities?: Opportunity[];
    /**
     * @type string | undefined
    */
    title?: string;
    /**
     * @type string | undefined
    */
    username?: string;
};