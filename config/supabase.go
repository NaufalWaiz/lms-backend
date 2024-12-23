package config

import "github.com/nedpals/supabase-go"

var SupabaseClient *supabase.Client

var SupabaseURL string = "https://vfhquqdbosumijamcqvo.supabase.co"

func InitSupabase() {
    url := "https://vfhquqdbosumijamcqvo.supabase.co" 
    key := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InZmaHF1cWRib3N1bWlqYW1jcXZvIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzI1OTMxNDcsImV4cCI6MjA0ODE2OTE0N30.eVULtSNNSuSih87mZJggUwfd61G47T8DPg9fFsbrt7M"

    SupabaseClient = supabase.CreateClient(url, key)
}