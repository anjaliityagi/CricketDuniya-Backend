BEGIN ;

CREATE TABLE otp_codes (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID REFERENCES users(id),
phone VARCHAR(20),
 otp_code VARCHAR(10) NOT NULL,
purpose VARCHAR(30) NOT NULL,
expires_at TIMESTAMP NOT NULL,
is_used BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP DEFAULT NOW()
);

COMMIT ;