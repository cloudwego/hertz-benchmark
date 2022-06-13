/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cpu

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRecordUsage(t *testing.T) {
	ctx, finished := context.WithCancel(context.Background())
	go func() {
		defer finished()
		sum := 0
		beginAt := time.Now()
		for {
			if time.Since(beginAt) > defaultInterval*5 {
				return
			}
			sum += 1
		}
	}()
	usage, err := RecordUsage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(usage.String())
}
