package config

//SiteSHortName ...
const SiteShortName string = "AdeWebGolang"

//SiteFullName ...
const SiteFullName string = "Adedayo Golang Web Application"

// SiteSlogan is widely use marketing words for the site.
const SiteSlogan string = "Hello Favs, Welcome to our Web Development Test App"

// SiteYear is the year the company starts it's operation.
const SiteYear int = 2021

// SiteRootTemplate is the root template folder location.
const SiteRootTemplate string = "html/"

// SiteDomainName define the full domain name of the site.
const SiteDomainName string = "adedayoadegboye.xyz"

// SiteProperDomainName define as a proper full domain name of the site.
const SiteProperDomainName string = "AdedayoAdegboye.xyz"

// SiteHeaderTemplate is the absolute path for the common header template for each HTML pages.
const SiteHeaderTemplate = SiteRootTemplate + "layout/header_front.html"

// SiteHeaderAccountTemplate is the absolute path for the common user account header template for each HTML pages.
const SiteHeaderAccountTemplate = SiteRootTemplate + "layout/header_account.html"

// SiteHeaderDashTemplate is the absolute path for the common dashboard header template for each HTML pages.
const SiteHeaderDashTemplate = SiteRootTemplate + "layout/header_dash.html"

// SiteFooterTemplate is the absolute path for the common footer template for each HTML pages.
const SiteFooterTemplate = SiteRootTemplate + "layout/footer_front.html"

// SiteFooterAccountTemplate is the absolute path for the common user account footer template for each HTML pages.
const SiteFooterAccountTemplate = SiteRootTemplate + "layout/footer_account.html"

// SiteFooterDashTemplate is the absolute path for the common dashboard footer template for each HTML pages.
const SiteFooterDashTemplate = SiteRootTemplate + "layout/footer_dash.html"

// SiteBaseURL is the base URL for the site URL structure. for prod https://adedayoadegboye.xyz/
// const SiteBaseURL = "http://127.0.0.1:8081/" //for dev
const SiteBaseURL = "https://adedayoadegboye.xyz/"

// const SiteBaseURL = "https://adedayoadegboye.xyz/"

// SiteTopMenuLogo is the small size top menu logo found at the top most left position.
const SiteTopMenuLogo = "/static/assets/images/adedayo_code_logo.png"

// EmailLogo is for email logo display on top of the email header content.
const EmailLogo = SiteBaseURL + "static/assets/images/adedayo_code_logoo.png"

// SiteEmail is the main technical support email for the company.
const SiteEmail = "majeedaadegboye@gmail.com"

// SitePhoneNumbers is the main contact numbers for the company.
const SitePhoneNumbers = ""

// SiteCompanyAddress is the company physical location.
const SiteCompanyAddress = "Your company address here"

// SiteTimeZone sets the default timezone to be used for this project.
const SiteTimeZone = "Asia/Manila"

// SecretKeyCORS is the secret key combination for the CORS (Cross-Origin Resource Sharing) middleware token.
const SecretKeyCORS = "n&@ix77r#^&^cgeb13w@!+pht^6qu-=("

// MyEncryptDecryptSK is for the Go's built-in encrypt and decrypt method.
const MyEncryptDecryptSK = "mkc&1*~#^8^#s0^=)^^7%a12"

// UserCookieExp is the user's cookie expiration in number of days.
const UserCookieExp = "30"

// SendGridAPIKey is the API key for the SendGrid SMTP server, make it encrypted later.
const SendGridAPIKey = "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"
