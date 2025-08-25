package tokens

type Token struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Interface string `json:"interface"`
		ID        string `json:"id"`
		Content   struct {
			Schema  string `json:"$schema"`
			JSONURI string `json:"json_uri"`
			Files   []struct {
				URI    string `json:"uri"`
				CdnURI string `json:"cdn_uri"`
				Mime   string `json:"mime"`
			} `json:"files"`
			Metadata struct {
				Description   string `json:"description"`
				Name          string `json:"name"`
				Symbol        string `json:"symbol"`
				TokenStandard string `json:"token_standard"`
			} `json:"metadata"`
			Links struct {
				Image string `json:"image"`
			} `json:"links"`
		} `json:"content"`
		Authorities []struct {
			Address string   `json:"address"`
			Scopes  []string `json:"scopes"`
		} `json:"authorities"`
		Compression struct {
			Eligible    bool   `json:"eligible"`
			Compressed  bool   `json:"compressed"`
			DataHash    string `json:"data_hash"`
			CreatorHash string `json:"creator_hash"`
			AssetHash   string `json:"asset_hash"`
			Tree        string `json:"tree"`
			Seq         int    `json:"seq"`
			LeafID      int    `json:"leaf_id"`
		} `json:"compression"`
		Grouping []interface{} `json:"grouping"`
		Royalty  struct {
			RoyaltyModel        string      `json:"royalty_model"`
			Target              interface{} `json:"target"`
			Percent             float64     `json:"percent"`
			BasisPoints         int         `json:"basis_points"`
			PrimarySaleHappened bool        `json:"primary_sale_happened"`
			Locked              bool        `json:"locked"`
		} `json:"royalty"`
		Creators []struct {
			Address  string `json:"address"`
			Share    int    `json:"share"`
			Verified bool   `json:"verified"`
		} `json:"creators"`
		Ownership struct {
			Frozen         bool        `json:"frozen"`
			Delegated      bool        `json:"delegated"`
			Delegate       interface{} `json:"delegate"`
			OwnershipModel string      `json:"ownership_model"`
			Owner          string      `json:"owner"`
		} `json:"ownership"`
		Supply    interface{} `json:"supply"`
		Mutable   bool        `json:"mutable"`
		Burnt     bool        `json:"burnt"`
		TokenInfo struct {
			Symbol       string `json:"symbol"`
			Supply       int64  `json:"supply"`
			Decimals     int    `json:"decimals"`
			TokenProgram string `json:"token_program"`
			PriceInfo    struct {
				PricePerToken float64 `json:"price_per_token"`
				Currency      string  `json:"currency"`
			} `json:"price_info"`
		} `json:"token_info"`
	} `json:"result"`
	ID string `json:"id"`
}
