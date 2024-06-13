import { Badge } from "./Badge";
import { Opportunity } from "./Opportunity";
import { UserProfile } from "./UserProfile";

 export type Collaboration = {
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
    hidden_at?: string;
    /**
     * @type integer | undefined
    */
    id?: number;
    /**
     * @type boolean | undefined
    */
    is_liked?: boolean;
    /**
     * @type boolean | undefined
    */
    is_payable?: boolean;
    /**
     * @type integer | undefined
    */
    likes_count?: number;
    opportunity?: Opportunity;
    /**
     * @type integer | undefined
    */
    opportunity_id?: number;
    /**
     * @type string | undefined
    */
    published_at?: string;
    /**
     * @type string | undefined
    */
    title?: string;
    user?: UserProfile;
    /**
     * @type integer | undefined
    */
    user_id?: number;
};