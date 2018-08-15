# errawr-gen

errawr-gen generates standard error constructors for errawr-compatible errors.
It reads a YAML file of errors and outputs code for a specified language.
Currently only Go is supported.

```
Usage of errawr-gen:
  -input-path string
      the path to read input from (default "-")
  -output-language string
      the language to write errors for (default "go")
  -output-path string
      the path to write output to (default "-")
  -package string
      the package to write
```

## Structure

### Domains

Error codes are constructed by combining a **domain**, a **section**, and an
**error**. Domains represent a unique project or namespace. They are identified
by a short abbreviation. For example:

* `idfc`: [insights-dataflow/controller](https://github.com/puppetlabs/insights-dataflow/tree/development/controller)
* `lidfp`: [insights-dataflow/proto](https://github.com/puppetlabs/insights-dataflow/tree/development/proto)

In general, error messages from complete applications should use an obvious
four-letter abbreviation. Other libraries and tools should can use a longer
abbreviation if necessary. (In the example above, `LIDFP` is "**l**ibrary for **i**nsights **d**ata**f**low **p**rotocol").

### Sections

A domain is divided into **sections**. A section represents a logically distinct
portion of the domain. Each domain will have its own rules and best practices
for defining sections, and some may have none at all.

For example, the insights-dataflow/controller domain might define the sections
`model` and `server`. Sections should usually be short identifiers constructed
using [snake case](https://en.wikipedia.org/wiki/Snake_case).

### Errors

Each section contains pertinent errors. The input file to the generator defines
the mapping of a domain to sections and sections to errors. For example:

```yaml
version: 1
domain:
  key: lsq
  title: Reflect SQL generation library
sections:
  driver:
    title: Driver errors
    errors:
      tcp_connection_error:
        title: TCP connection error
        description:
          friendly: >
            We could not access this service. You may need to check your
            firewall configuration.
          technical: >
            The host {{code host}} is not connectable on TCP port {{port}}.
        arguments:
          host:
            description: the host name
          port:
            type: integer
            validators:
              - positive_number
              - integer
            description: the TCP port number
      authentication_error:
        title: Authentication error
        description: >
          We could not authenticate to this service with the credentials
          provided.
```

All errors must have, at minimum, a `title` and a `description`.

#### Descriptions

Frequently, descriptions are just text strings.

A `description` may optionally be an object with one or more of the following
keys: `friendly`, `technical`. Descriptions are
[Handlebars](http://handlebarsjs.com/) templates and substitute variables from
the `arguments` mapping. Additionally, several formatting helpers are available,
roughly following the style of HTML:

| Helper | HTML equivalent | Text equivalent |
| ------ | --------------- | --------------- |
| `{{em text}}` | `<em>{{text}}</em>` | ``*{{text}}*`` |
| `{{link url text}}` | `<a href="{{url}}">{{text}}</a>` | ``{{text}} ({{url}})`` |
| `{{#join list}}#{{@index}}{{/join}}` | `<ul><li>#0</li><li>#1</li></ul>` | ``#1 and #2`` |
| `{{pre text}}` | `<code>{{text}}</code>` | `` `{{text}}` `` |
| `{{quote text}}` | `&ldquo;{{text}}&rdquo;` | `"{{text}}"` |

#### Arguments

Arguments are mapping of a name (which can be substituted in the description
template as described above) to an argument definition. No information is
required to be included in an argument definition; mapping entries may simply be
empty values.

A `description` of an argument may be provided, which can be useful for
supporting users when trying to understand an error. Generally, the description
is not surfaced to end users directly.

A `type` may be provided. Types mostly map directly to JSON, but are a bit more
constrained. Valid types:

* `string`
* `number`
* `integer`
* `boolean`
* `list<string>`
* `list<number>`
* `list<integer>`
* `list<boolean>`

A `default` may be provided, which will be substituted in the description
template if no value is provided for the argument when the error is created. If
this key is not defined, a value for the argument must be supplied when the
error is created.

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

An **error code** is the concatenation of the code's domain, section, and name,
separated by underscores. In the example above, the error code is
`lsq_driver_tcp_connection_error`.

These representations are unique within the entire error collection and can be
used by support team members to easily identify a specific error.

### Errors

When rendered, an error has the following object structure:

```json
{
  "code": "lsq_driver_tcp_connection_error",
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
descriptions, use the pronoun "we" to refer to Puppet as a collective, and the
pronoun "you" to refer to the end user when they can take action on their own.

Never use gendered pronouns in any part of error definitions.

Argument descriptions should be simple phrases like "the host name." Complete
sentences are not necessary.

## Libraries

* Go: [errawr-go](https://github.com/puppetlabs/errawr-go)
