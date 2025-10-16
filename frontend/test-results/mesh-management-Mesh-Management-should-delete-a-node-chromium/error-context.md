# Page snapshot

```yaml
- generic [ref=e4]:
  - generic [ref=e5]:
    - heading "Create your account" [level=2] [ref=e6]
    - paragraph [ref=e7]:
      - text: Or
      - link "sign in to your existing account" [ref=e8] [cursor=pointer]:
        - /url: /login
  - generic [ref=e9]:
    - paragraph [ref=e11]: Email already exists
    - generic [ref=e12]:
      - generic [ref=e13]:
        - generic [ref=e14]: Display name
        - textbox "Display name" [ref=e15]: Test User
      - generic [ref=e16]:
        - generic [ref=e17]: Email address
        - textbox "Email address" [ref=e18]: test-1760598163760@example.com
      - generic [ref=e19]:
        - generic [ref=e20]: Password
        - textbox "Password" [ref=e21]: testpassword123
    - button "Create account" [ref=e23]
```