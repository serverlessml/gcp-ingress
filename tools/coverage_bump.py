#! /usr/local/bin/python3

# Copyright 2020 dkisler.com Dmitry Kisler
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
# OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND
# NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS
# BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN
# ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF, OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# See the License for the specific language governing permissions and
# limitations under the License.

"""Tool to patch the README.md with the code coderage badde

https://shields.io/ is used to generate badges
"""

import argparse
import logging
import re
from pathlib import Path

import utils  # type: ignore


def get_args() -> argparse.Namespace:
    """Parse input arguments."""
    parser = argparse.ArgumentParser("Cover test coverage assessment.")
    parser.add_argument(
        "--platform", "-p", help="Env platform, i.e. aws, or gcp.", required=True, type=str
    )
    args = parser.parse_args()
    return args


def run_gocover(path: Path, platform: str) -> None:
    """Run gocover."""
    cmd = f"""export TEMP_DIR=/tmp/temp-ingress \
&& mkdir -p $TEMP_DIR/bus \
&& cp -r handlers config $TEMP_DIR \
&& cp bus/{platform}*.go $TEMP_DIR/bus/ \
&& cp main_{platform}.go runner.go go.* $TEMP_DIR/ \
&& cd $TEMP_DIR \
&& go mod tidy \
&& go test -tags test -coverprofile=$TEMP_DIR/go-cover.tmp ./... > /dev/null \
&& go tool cover -func $TEMP_DIR/go-cover.tmp -o {path} \
&& cd /tmp && rm -r $TEMP_DIR"""

    utils.execute_cmd(cmd)


def extract_total_coverage(raw: str) -> int:
    """Extract total coverage."""
    tail_line = raw.splitlines()[-1]
    return int(float(tail_line.split("\t")[-1][:-1]))


def generate_url(coverage_pct: float, platform: str) -> str:
    """Generate badge source URL."""
    color = "yellow"
    if coverage_pct == 100:
        color = "brightgreen"
    elif coverage_pct > 90:
        color = "green"
    elif coverage_pct > 70:
        color = "yellowgreen"
    elif coverage_pct > 50:
        color = "yellow"
    else:
        color = "orange"

    return f"https://img.shields.io/badge/coverage%20{platform}-{coverage_pct}%25-{color}"


def main(platform: str) -> None:
    """
    Run.

    Args:
        platform: Infra env AWS, or GCP.
    """
    root = Path(__file__).absolute().parents[1]
    path_readme = root / "README.md"
    path_coverage = root / f"COVERAGE.{platform}"
    placeholder_tag = f"code coverage:{platform}"
    regexp_pattern = rf"\[\!\[{placeholder_tag}\]\(.*\)\]\(.*\)"

    run_gocover(path_coverage, platform)

    coverage = utils.read(path_coverage)

    coverage_pct = extract_total_coverage(coverage)

    badge_url = generate_url(coverage_pct, platform)

    inpt = utils.read(path_readme)

    search = re.findall(regexp_pattern, inpt)

    if not search:
        raise Exception(f"No placeholder found in README.md. Add '[![{placeholder_tag}]()]()'.")

    placeholder_inject = f"[![{placeholder_tag}]({badge_url})]({badge_url})"

    out = re.sub(regexp_pattern, placeholder_inject, inpt)

    utils.write(out, path_readme)


if __name__ == "__main__":
    log = logging.getLogger("coverage-bump")

    args = get_args()

    try:
        main(args.platform)
    except Exception as ex:  # pylint: disable=broad-except
        log.error(ex)
