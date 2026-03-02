CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,        
    user_id UUID NOT NULL,       
    access_id UUID NOT NULL,      
    exp_at TIMESTAMP WITH TIME ZONE NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP 
);


CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);