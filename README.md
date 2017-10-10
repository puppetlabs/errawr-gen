# reflect-errors

This is a canonical collection of errors used by Reflect applications. The error
messages in this repository are designed to be surfaced to end users; they
should not be used only for internal applications.

This repository does not contain any code. Instead, other repositories should
reference the contents of this repository as a source of truth.

## Structure

### Domains

Error codes are constructed by combining a **domain**, a **section**, and an
**error**. Domains represent a unique part of the Reflect system. They are
identified by a short abbreviation. For example:

* `RA`: [reflect-api](https://github.com/reflect/reflect-api)
* `RP`: [reflect-app](https://github.com/reflect/reflect-app)
* `RG`: [reflect-agent](https://github.com/reflect/reflect-agent)
* `LRE`: [reflect-reporting](https://github.com/reflect/reflect-reporting)
* `LSQ`: [reflect-sql](https://github.com/reflect/reflect-sql)

In general, error messages from complete applications should use a two-letter
domain abbreviation starting with `R` (for Reflect, of course). Other libraries
and tools should use a sensible three-letter abbreviation. (In the example
above, `LRE` could be read as "**l**ibrary, **re**porting.")

Domains are configured in the `domains.yml` file in the `domains` directory.
Domain-specific information is contained in a directory that matches the `name`
key of each domain defined.

Domains can be allocated in advance. If no directory exists to match the name of
the domain, it is reserved but not processed.

### Sections

A domain is divided into **sections**. A section is a two-digit identifier that
is used to represent a logically distinct portion of the domain. Each domain
will have its own rules and best practices for defining sections, and some may
have none at all.

Section numbers begin at `10` and continue through `99`. A section number
starting with `0` is not valid. In a given domain directory, sections are
configured in the file `sections.yml`. Section-specific information is contained
in a file that matches the `name` key of each section defined, prefixed with an
underscore.

For example, the reflect-api domain has a section for authentication errors. In
`reflect-api/sections.yml`, this section is given the number `14` and the name
`authentication`. Configuration for the section is located in the file
`reflect-api/_authentication.yml`.

Like domains, sections can be allocated in advance. If no section configuration
file exists, no errors are defined for that section.

### Errors

Each section contains **errors** that are relevant to that section. Errors are
uniquely identified by three-digit numbers, which begin at `100` and continue
through `999`. An error identifier starting with `0` is not valid.

Section files contain mappings of error numbers to error definitions. For
example:

```yaml
errors:
  101:
    title: TCP connection error
    description:
      friendly: |
        We could not access this service. You may need to check your firewall
        configuration.
      technical: |
        The host {{code host}} is not connectable on TCP port {{port}}.
    arguments:
      host:
        description: the host name
      port:
        validators:
          - positive_number
          - integer
        description: the TCP port number
  102:
    title: Authentication error
    description: |
      We could not authenticate to this service with the credentials provided.
```

All errors must have, at minimum, a `title` and a `description`.

#### Descriptions

Frequently, descriptions are just text strings.

A `description` may optionally be an object with one or more of the following
keys: `friendly`, `technical`. Descriptions are
[Handlebars](http://handlebarsjs.com/) templates and substitute variables from
the `arguments` mapping. Additionally, several formatting helpers are available,
roughly following the style of HTML:

| Helper | HTML equivalent |
| ------ | --------------- |
| `{{code text}}` | `<code>{{text}}</code>` |
| `{{em text}}` | `<em>{{text}}</em>` |
| `{{link url text}}` | `<a href="{{url}}">{{text}}</a>` |

#### Arguments

Arguments are mapping of a name (which can be substituted in the description
template as described above) to an argument definition. No information is
required to be included in an argument definition; mapping entries may simply be
empty values.

A `description` of an argument may be provided, which can be useful for
supporting users when trying to understand an error. Generally, the description
is not surfaced to end users directly.

A `default` may be provided, which will be substituted in the description
template if no value is provided for the argument when the error is created. If
this key is not defined, the default value is `null`.

A list of `validators` may be provided. Validators check that the arguments are
sane before allowing them to be surfaced and used in the description templates.
If any validator fails, the argument is invalidated, and the default value is
used instead. The following validators are available:

| Name | Description |
| ---- | ----------- |
| `number` | Asserts that the argument is a decimal number. |
| `positive_number` | Like `number`, but also checks that the argument is greater than 0. |
| `nonnegative_number` | Like `number`, but also checks that the argument is at least 0. |
| `integer` | Asserts that the argument is an integer. |

## Representation

### Error codes

An **error code** is the concatenation of the code's domain, section, and
numerical identifier. For example, code `162` of the reflect-api authentication
domain would be represented as `RA14162` to an end user.

These representations are unique within the entire error collection and can be
used by support team members to easily identify a specific error.

### Errors

When rendered, an error has the following object structure:

```json
{
  "code": "LSQ10101",
  "title": "TCP connection error",
  "description": {
    "friendly": "We could not access this service. You may need to check your firewall configuration.",
    "technical": "The host {{code host}} is not connectable on TCP port {{port}}."
  },
  "arguments": {
    "host": "localhost",
    "port": "5432"
  },
  "formatted": {
    "friendly": "We could not access this service. You may need to check your firewall configuration.",
    "technical": "The host `localhost` is not connectable on TCP port 5432."
  },
  "causes": []
}
```

These property semantics should be fairly intuitive. Of note, however, is the
`formatted` property, in which the error renderer has substituted the arguments
into the description templates, and the `causes` property, which contains an
optional list of more error objects that caused the given error to occur. If
there are no additional causes, the `causes` property may be omitted altogether.

## Style

In general, error descriptions should always be complete sentences. They should
always end with a punctuation mark, generally a period.

Use of the passive voice ("the account is deactivated") should be avoided where
possible. It may occur in technical descriptions, but should be used sparingly
in friendly descriptions as users tend to find it inaccessible. In friendly
descriptions, use the pronoun "we" to refer to Reflect as a collective, and the
pronoun "you" to refer to the end user when they can take action on their own.

Never use gendered pronouns in any part of error definitions.

Argument descriptions should be simple phrases like "the host name." Complete
sentences are not necessary.

## Libraries

* Go: [reflect-errors-go](https://github.com/reflect/reflect-errors-go)
